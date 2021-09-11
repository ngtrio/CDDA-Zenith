package data

import "github.com/tidwall/gjson"

type Mod struct {
	ID           string
	Name         string
	Description  string
	Path         string
	Dependencies []string
	IdMap        map[string][]string
	NameMap      map[string][]string
	TempData     map[string][]*gjson.Result
	Loaded       bool
}

func (mod *Mod) GetById(id string) []string {
	if v, ok := mod.IdMap[id]; ok {
		return v
	}
	return nil
}
