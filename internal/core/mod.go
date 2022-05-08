package core

import (
	"github.com/tidwall/gjson"
	"zenith/internal/loader"
)

type Mod struct {
	Id           string
	Name         string
	Description  string
	Path         string
	Dependencies []string
	TempData     map[string]map[string]*gjson.Result
	Loaded       bool
}

func (mod *Mod) CreateIndex(indexer Indexer, json *gjson.Result, langPack map[string]LangPack) {

	if !loader.IsAbstract(json) {
		for _, pack := range langPack {
			vo := NewVO(mod.Id, mod.Name)
			vo.Bind(json, pack, mod, indexer)
			indexer.AddRangeIndex(vo)
			indexer.AddNameIndex(vo)
			indexer.AddIdIndex(vo)
		}
	}
}
