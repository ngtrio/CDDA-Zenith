package core

import (
	"fmt"
	"strconv"
	"strings"
	"zenith/internal/constdef"
	"zenith/internal/loader"
	"zenith/internal/view"
	"zenith/pkg/fileutil"
	"zenith/pkg/jsonutil"

	"github.com/leonelquinteros/gotext"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type Game struct {
	Commit    string
	UpdateAt  string
	Mods      map[string]*Mod // id -> mod
	ModPath   string
	Lang      string
	Mo        *gotext.Mo
	Po        *gotext.Po
	TypeItems map[string][]*gjson.Result // map(type -> items)
}

func (game *Game) Load(targets map[string]bool) {
	if err := game.preLoad(); err != nil {
		log.Fatal(err)
	}

	for _, mod := range game.Mods {
		if len(targets) > 0 {
			if _, ok := targets[mod.ID]; ok && !mod.Loaded {
				game.doLoad(mod)
			}
		} else {
			game.doLoad(mod)
		}
	}

	game.postLoad()
}

func (game *Game) preLoad() error {

	if _, dirs, err := fileutil.Ls(game.ModPath); err != nil {
		return err
	} else {
		for _, dir := range dirs {
			if files, _, err := fileutil.Ls(dir); err != nil {
				return err
			} else {
				for _, file := range files {
					if strings.HasSuffix(file, "modinfo.json") {
						res := loader.LoadJsonFromFile(file)
						modInfo := *res[0]
						id := modInfo.Get("id").String()

						var path string
						jsonPath := modInfo.Get("path")
						if !jsonPath.Exists() {
							path = dir
						} else {
							path = dir + "/" + jsonPath.String()
						}
						var dependencies []string
						jsonDp := modInfo.Get("dependencies")
						if jsonDp.Exists() {
							for _, dp := range jsonDp.Array() {
								dependencies = append(dependencies, dp.String())
							}
						}

						mod := &Mod{
							ID:           id,
							Name:         modInfo.Get("name").String(),
							Description:  modInfo.Get("description").String(),
							Path:         path,
							Dependencies: dependencies,
							IdMap:        make(map[string][]*gjson.Result),
							NameMap:      make(map[string][]*gjson.Result),
							TypeItems:    make(map[string][]*gjson.Result),
							TempData:     make(map[string][]*gjson.Result),
							Loaded:       false,
						}

						game.Mods[id] = mod
					}
				}
			}
		}
	}

	return nil
}

func (game *Game) doLoad(mod *Mod) {
	if mod.Loaded {
		return
	}

	dependencies := mod.Dependencies
	for _, dependency := range dependencies {
		m := game.Mods[dependency]
		if m != nil {
			game.doLoad(m)
		} else {
			log.Warnf("%v's dependency: %v is not found.", mod.Name, dependency)
		}
	}
	path := mod.Path
	jsons := loader.LoadJsonFromPaths(path)
	game.processModData(mod, jsons)

	mod.Loaded = true
}

func (game *Game) postLoad() {
	game.Mo = loader.LoadMo(game.Lang)
	game.Po = loader.LoadPo(game.Lang)

	for _, mod := range game.Mods {
		mod.Finalize(game.Mo)
		for t, i := range mod.TypeItems {
			game.TypeItems[t] = append(game.TypeItems[t], i...)
		}
	}
}

func (game *Game) processModData(mod *Mod, jsons []*gjson.Result) {
	for _, json := range jsons {
		id := getId(json)
		if id == "" {
			log.Debugf("id not found, json: %s", json.String())
			continue
		}

		if !isInAllowList(json) {
			continue
		}

		tar := mod.TempData[id]
		mod.TempData[id] = append(tar, json)
	}

	for _, tempJsonList := range mod.TempData {
		for _, tempJson := range tempJsonList {
			if loader.NeedInherit(tempJson) {
				game.inherit(mod, tempJson)
			}
		}
	}
}

func (game *Game) inherit(mod *Mod, json *gjson.Result) bool {
	cf := json.Get("copy-from")
	if !cf.Exists() {
		return false
	}
	parId := cf.String()

	flag := false
	if pars := mod.TempData[parId]; pars != nil {
		for _, par := range pars {
			if par != json && par.Get("type").String() == json.Get("type").String() {
				if loader.NeedInherit(par) {
					game.inherit(mod, par)
				}

				jsonStr := par.String()
				json.ForEach(func(k, v gjson.Result) bool {
					field := k.String()
					switch field {
					case "relative":
						inheritRelative(&jsonStr, par, json, "relative")
					case "proportional":
						inheritProportional(&jsonStr, par, json, "proportional")
					case "extend":
						v.ForEach(func(ck, cv gjson.Result) bool {
							vInPar := par.Get(ck.String())
							var res []interface{}
							if vInPar.Exists() {
								for _, elem := range vInPar.Array() {
									res = append(res, elem.Value())
								}
							}
							for _, elem := range cv.Array() {
								res = append(res, elem.Value())
							}

							jsonutil.Set(&jsonStr, ck.String(), res)
							return true
						})
					case "delete":
						v.ForEach(func(ck, cv gjson.Result) bool {
							vInCur := gjson.Get(jsonStr, ck.String())
							if vInCur.Exists() {
								var res []string

								// FIXME fully support delete
								if !vInCur.IsArray() {
									id := json.Get("id").String()
									log.Warnf("delete field is not supported, id: %v", id)
									return true
								}

								for _, elem := range vInCur.Array() {
									flag := false
									for _, cvElem := range cv.Array() {
										if elem.String() == cvElem.String() {
											flag = true
										}
										break
									}
									if !flag {
										res = append(res, elem.String())
									}
								}
								jsonutil.Set(&jsonStr, ck.String(), res)
							}

							// we assume that delete is done from self
							if par.Get(ck.String()).Exists() && !vInCur.Exists() {
								log.Debugf("%s field delete is abnormal", json)
							}
							return true
						})

					case "copy-from":
						// discard
					default:
						jsonutil.Set(&jsonStr, k.String(), v.Value())
					}
					return true
				})
				*json = gjson.Parse(jsonStr)
				flag = true

				break
			}
		}
	}

	if !flag {
		for _, dp := range mod.Dependencies {
			dpMod := game.Mods[dp]
			if game.inherit(dpMod, json) {
				break
			}
		}
	}

	return true
}

func getId(json *gjson.Result) string {
	var id string
	var has bool
	if id, has = jsonutil.GetString("id", json, ""); has {
		return id
	}
	if id, has = jsonutil.GetString("abstract", json, ""); has {
		return id
	}
	return ""
}

func inheritRelative(jsonStr *string, par *gjson.Result, json *gjson.Result, path string) {
	json.Get(path).ForEach(func(k, v gjson.Result) bool {

		if v.IsObject() {
			inheritRelative(jsonStr, par, json, path+"."+k.String())
		} else if v.Type == gjson.Number {
			fromPath := path + "." + k.String()
			toPath := strings.Split(fromPath, "relative.")[1]
			parVal := par.Get(toPath)
			if parVal.Exists() {
				jsonutil.Set(jsonStr, toPath, parVal.Int()+v.Int())
			}
		}
		return true
	})
}

func inheritProportional(jsonStr *string, par *gjson.Result, json *gjson.Result, path string) {
	json.Get(path).ForEach(func(k, v gjson.Result) bool {

		if v.IsObject() {
			inheritProportional(jsonStr, par, json, path+"."+k.String())
		} else if v.Type == gjson.Number {
			fromPath := path + "." + k.String()

			toPath := strings.Split(fromPath, "proportional.")[1]

			parVal := par.Get(toPath)

			if parVal.Exists() {
				val, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", parVal.Float()*v.Float()), 64)
				jsonutil.Set(jsonStr, toPath, val)
			}
		}
		return true
	})
}

func isInAllowList(json *gjson.Result) bool {
	type_, _ := jsonutil.GetString("type", json, "")

	// TODO add new type here
	allowList := map[string]bool{"MONSTER": true}
	return allowList[type_]
}

func (game *Game) GetById(id string) []*view.VO {
	return game.GetByModAndId("", id)
}

func (game *Game) GetByModAndId(modId, id string) []*view.VO {
	jsons := make(map[string][]*gjson.Result)
	if modId == "" {
		for _, mod := range game.Mods {
			jsons[mod.ID] = append(jsons[mod.ID], mod.GetById(id)...)
		}
	} else {
		mod := game.Mods[modId]
		jsons[modId] = append(jsons[modId], mod.GetById(id)...)
	}
	return game.jsonToVO(jsons)
}

func (game *Game) GetByName(name string) []*view.VO {
	return game.GetByModAndName("", name)
}

func (game *Game) GetByModAndName(modId, name string) []*view.VO {
	jsons := make(map[string][]*gjson.Result)
	if modId == "" {
		for _, mod := range game.Mods {
			jsons[mod.ID] = append(jsons[mod.ID], mod.GetByName(name)...)
		}
	} else {
		mod := game.Mods[modId]
		jsons[modId] = append(jsons[modId], mod.GetByName(name)...)
	}
	return game.jsonToVO(jsons)
}

func (game *Game) GetByType(type_ string, num, page int) ([]*view.VO, int) {
	return game.GetByModAndType("", type_, num, page)
}

func (game *Game) GetByModAndType(modId, type_ string, num, page int) (res []*view.VO, totalPage int) {
	jsons := make(map[string][]*gjson.Result)
	if modId == "" {
		if arr, has := game.TypeItems[type_]; !has {
			return
		} else {
			total := len(arr)
			totalPage = total / num
			if total%num > 0 {
				totalPage += 1
			}
			start := num * page
			end := start + num - 1
			if start >= total {
				return
			}
			if end > total {
				end = total
			}

			for i := start; i < end; i++ {
				modId, _ = jsonutil.GetString(constdef.FieldModId, arr[i], "")
				jsons[modId] = append(jsons[modId], arr[i])
			}
		}
	} else {
		if m, has := game.Mods[modId]; has {
			jsons[modId], totalPage = m.GetByType(type_, num, page)
		}
	}

	res = game.jsonToVO(jsons)
	return
}

func (game *Game) FuzzyGet(keyword, tp string) []*view.VO {
	if len(keyword) <= 0 {
		return nil
	}

	res := make(map[string][]*gjson.Result, 0)
	for _, mod := range game.Mods {
		for name, items := range mod.NameMap {
			if strings.Contains(name, keyword) {
				res[mod.ID] = append(res[mod.ID], items...)
			}
		}
	}

	return game.jsonToVO(res)
}

func (game *Game) jsonToVO(jsonMap map[string][]*gjson.Result) []*view.VO {
	vos := make([]*view.VO, 0)
	for modId, jsons := range jsonMap {
		mod := game.Mods[modId]
		for _, json := range jsons {
			tp, _ := jsonutil.GetString(constdef.FieldType, json, "")
			vo := view.NewVO(mod.Name)
			vo.Bind(tp, json, game.Mo, game.Po)
			vos = append(vos, vo)
		}
	}

	return vos
}
