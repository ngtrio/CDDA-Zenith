package data

import (
	"zenith/internal/i18n"
	"zenith/internal/loader"

	"github.com/leonelquinteros/gotext"
	"github.com/tidwall/gjson"
)

type Mod struct {
	ID           string
	Name         string
	Description  string
	Path         string
	Dependencies []string
	IdMap        map[string][]*string
	NameMap      map[string][]*string
	TempData     map[string][]*gjson.Result
	Loaded       bool
}

func (mod *Mod) GetById(id string) []string {
	return ptrToValue(mod.IdMap[id])
}

func (mod *Mod) GetByName(name string) []string {
	return ptrToValue(mod.NameMap[name])
}

func ptrToValue(ptrs []*string) []string {
	res := make([]string, 0)
	for _, v := range ptrs {
		res = append(res, *v)
	}
	return res
}

func (mod *Mod) Finalize(mo *gotext.Mo) {
	for id, jsons := range mod.TempData {
		for _, json := range jsons {
			if !loader.IsAbstract(json) {
				mod.processType(json)
				mod.createIndex(id, json, mo)
			}
		}
	}
}

func (mod *Mod) processType(json *gjson.Result) {

}

func (mod *Mod) createIndex(id string, json *gjson.Result, mo *gotext.Mo) {
	jsonStr := json.String()
	name := i18n.Tran("name", json, mo)
	mod.IdMap[id] = append(mod.IdMap[id], &jsonStr)
	mod.NameMap[name] = append(mod.NameMap[name], &jsonStr)
}
