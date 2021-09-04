package core

import (
	"strings"
	"zenith/internal/data"
	"zenith/internal/loader"
	"zenith/pkg/fileutil"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type Game struct {
	Version string
	Mods    map[string]*data.Mod
	ModPath string
}

func (game *Game) LoadMod(targets map[string]bool) {

	if err := game.preLoadMod(); err != nil {
		log.Fatal(err)
	}

	for _, mod := range game.Mods {
		if len(targets) > 0 {
			if _, ok := targets[mod.ID]; ok && !mod.Loaded {
				game.doLoadMod(mod)
			}
		} else {
			game.doLoadMod(mod)
		}
	}
}

func (game *Game) preLoadMod() error {

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

						mod := &data.Mod{
							ID:           id,
							Name:         modInfo.Get("name").String(),
							Description:  modInfo.Get("description").String(),
							Path:         path,
							Dependencies: dependencies,
							Data:         make(map[string][]string),
							TempData:     make(map[string][]string),
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

func (game *Game) doLoadMod(mod *data.Mod) {
	dependencies := mod.Dependencies
	for _, dependency := range dependencies {
		m := game.Mods[dependency]
		if !m.Loaded {
			game.doLoadMod(m)
		}
	}
	path := mod.Path
	// TODO
	jsons := loader.LoadJsonFromPaths(path)
	game.processModData(mod, jsons)

	mod.Loaded = true
	log.Printf("[MOD]: %s is loaded", mod.ID)
}

func (game *Game) processModData(mod *data.Mod, jsons []*gjson.Result) {
	for _, json := range jsons {
		id, isAbstract := getId(json)
		if id == "" {
			log.Debug("id not found, json: %s", json.String())
			continue
		}

		var res string
		if needInherit(json) {
			res = game.inherit(json)
		} else {
			res = json.String()
		}

		if isAbstract {
			tar := mod.TempData[id]
			mod.TempData[id] = append(tar, res)
		} else {
			tar := mod.Data[id]
			mod.Data[id] = append(tar, res)
		}
	}
}

func needInherit(json *gjson.Result) bool {
	return json.Get("copy-from").Exists()
}

func getId(json *gjson.Result) (string, bool) {
	id := json.Get("id")
	if id.Exists() {
		return id.String(), false
	}
	ab := json.Get("abstract")
	if ab.Exists() {
		return ab.String(), true
	}
	return "", false
}

func (game *Game) inherit(json *gjson.Result) string {
	cf := json.Get("copy-from")
	if cf.Exists() {
		// parId := cf.String()
	}
	return ""
}
