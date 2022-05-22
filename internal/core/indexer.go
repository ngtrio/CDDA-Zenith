package core

import (
	log "github.com/sirupsen/logrus"
	"strings"
	"zenith/internal/constdef"
)

type rangeIndex map[string][]*VO
type nameIndex map[string]map[string][]*VO
type idIndex map[string]map[string][]*VO
type modIndexGroup struct {
	rangeIndex  rangeIndex
	nameIndexes nameIndex // lang -> nameIndex
	idIndex     idIndex
}
type allModIndexGroup map[string]modIndexGroup
type i18nIndexGroup map[string]allModIndexGroup

type Indexer interface {
	RangeIndex(tp string, num, page int, lang string) ([]*VO, int)
	ModRangeIndex(mod, tp string, num, page int, lang string) ([]*VO, int)
	NameIndex(tp, name, lang string) []*VO
	ModNameIndex(mod, tp, name, lang string) []*VO
	FuzzyNameIndex(tp, keyword, lang string) []*VO
	ModFuzzyNameIndex(mod, tp, keyword, lang string) []*VO
	IdIndex(tp, id, lang string) []*VO
	ModIdIndex(mod, tp, id, lang string) []*VO
	TypeIndex(tp, lang string) []*VO

	AddRangeIndex(vo *VO)
	AddNameIndex(vo *VO)
	AddIdIndex(vo *VO)

	PrintItemNum()
}

func NewMemIndexer(mods map[string]*Mod, langPacks map[string]LangPack) Indexer {
	indexer := &MemIndexer{
		BaseIndexer: BaseIndexer{
			i18nIndexes: make(i18nIndexGroup),
		},
	}

	for _, mod := range mods {
		indexer.modIds = append(indexer.modIds, mod.Id)
		for lang := range langPacks {
			if _, has := indexer.i18nIndexes[lang]; !has {
				indexer.i18nIndexes[lang] = make(allModIndexGroup)
			}
			indexer.i18nIndexes[lang][mod.Id] = modIndexGroup{
				rangeIndex:  make(rangeIndex),
				nameIndexes: make(nameIndex),
				idIndex:     make(idIndex),
			}
		}
	}

	return indexer
}

type BaseIndexer struct {
	i18nIndexes i18nIndexGroup
	modIds      []string // keep mod order
}

type MemIndexer struct {
	BaseIndexer
}

func (i *MemIndexer) PrintItemNum() {
	for lang, idxes := range i.i18nIndexes {
		for mod, idx := range idxes {
			for tp, items := range idx.idIndex {
				num := 0
				for _, its := range items {
					num += len(its)
				}
				log.Infof("lang: %v, mod: %v, type: %v, num: %v \n", lang, mod, tp, num)
			}
		}
	}
}

func (i *MemIndexer) RangeIndex(tp string, num, page int, lang string) ([]*VO, int) {
	if _, has := i.i18nIndexes[lang]; !has {
		return nil, 0
	}

	total := 0
	for _, idxes := range i.i18nIndexes[lang] {
		total += len(idxes.rangeIndex[tp])
	}
	totalPage := total / num
	if total%num > 0 {
		totalPage++
	}

	if page < 0 {
		return nil, totalPage
	}
	res := make([]*VO, 0, num)
	offset := num * (page - 1)
	for _, modId := range i.modIds {
		if idx, has := i.i18nIndexes[lang][modId]; has {
			l := len(idx.rangeIndex[tp])

			if offset >= l {
				offset -= l
				continue
			}

			if num <= 0 {
				break
			}

			if num+offset <= l {
				res = append(res, idx.rangeIndex[tp][offset:offset+num]...)
				break
			} else {
				res = append(res, idx.rangeIndex[tp][offset:]...)
				num = num - (l - offset)
			}
		}
	}

	return res, totalPage
}

func (i *MemIndexer) ModRangeIndex(mod, tp string, num, page int, lang string) ([]*VO, int) {
	if _, has := i.i18nIndexes[lang]; !has {
		return nil, 0
	}

	if idxes, has := i.i18nIndexes[lang][mod]; !has || len(idxes.rangeIndex[tp]) == 0 {
		return nil, 0
	} else {
		total := len(idxes.rangeIndex[tp])
		totalPage := total / page
		if total%page > 0 {
			totalPage++
		}

		offset := num * (page - 1)
		end := offset + num
		if end >= len(idxes.rangeIndex) {
			end = len(idxes.rangeIndex) - 1
		}
		return idxes.rangeIndex[tp][offset:end], totalPage
	}
}

func (i *MemIndexer) NameIndex(tp, name, lang string) []*VO {
	if _, has := i.i18nIndexes[lang]; !has {
		return nil
	}

	res := make([]*VO, 0)
	for _, modId := range i.modIds {
		res = append(res, i.ModNameIndex(modId, tp, name, lang)...)
	}

	return res
}

func (i *MemIndexer) ModNameIndex(mod, tp, name, lang string) []*VO {
	if _, has := i.i18nIndexes[lang]; !has {
		return nil
	}

	res := make([]*VO, 0)
	if idx, has := i.i18nIndexes[lang][mod]; has {
		if nameMap, has := idx.nameIndexes[tp]; has {
			if items, has := nameMap[name]; has {
				res = append(res, items...)
			}
		}
	}

	return res
}

func (i *MemIndexer) FuzzyNameIndex(tp, keyword, lang string) []*VO {
	if _, has := i.i18nIndexes[lang]; !has {
		return nil
	}

	temp := make([]*VO, 0)
	for _, modId := range i.modIds {
		temp = append(temp, i.ModFuzzyNameIndex(modId, tp, keyword, lang)...)
	}

	exact := make([]*VO, 0)
	left := make([]*VO, 0)
	for _, vo := range temp {
		if vo.Name == keyword {
			exact = append(exact, vo)
		} else {
			left = append(left, vo)
		}
	}

	return append(exact, left...)
}

func (i *MemIndexer) ModFuzzyNameIndex(mod, tp, keyword, lang string) []*VO {
	if _, has := i.i18nIndexes[lang]; !has {
		return nil
	}

	res := make([]*VO, 0)
	if idx, has := i.i18nIndexes[lang][mod]; has {
		if nameMap, has := idx.nameIndexes[tp]; has {
			for name, items := range nameMap {
				if strings.Contains(name, keyword) {
					res = append(res, items...)
				}
			}
		}
	}

	return res
}

func (i *MemIndexer) IdIndex(tp, id, lang string) []*VO {
	var tps []string
	if tp == constdef.TypeItem {
		for t := range constdef.ItemTypes {
			tps = append(tps, t)
		}
	} else {
		tps = []string{tp}
	}

	res := make([]*VO, 0)
	for _, tp := range tps {
		for _, modId := range i.modIds {
			r := i.ModIdIndex(modId, tp, id, lang)
			if r != nil {
				res = append(res, r...)
			}
		}
	}
	return res
}

func (i *MemIndexer) ModIdIndex(mod, tp, id, lang string) []*VO {
	if _, has := i.i18nIndexes[lang]; !has {
		return nil
	}

	res := make([]*VO, 0)
	if idx, has := i.i18nIndexes[lang][mod]; has {
		if idMap, has := idx.idIndex[tp]; has {
			if items, has := idMap[id]; has {
				res = append(res, items...)
			}
		}
	}

	return res
}

func (i *MemIndexer) TypeIndex(tp, lang string) []*VO {
	res := make([]*VO, 0)
	if modIdx, has := i.i18nIndexes[lang]; has {
		for _, idx := range modIdx {
			for _, vo := range idx.idIndex[tp] {
				res = append(res, vo...)
			}
		}
	}

	return res
}

func (i *MemIndexer) AddRangeIndex(vo *VO) {
	mod := vo.ModId
	tp := vo.Type
	lang := vo.Lang

	if _, has := i.i18nIndexes[lang]; !has {
		log.Warnf("lang: %v not found", lang)
		return
	}

	if _, has := i.i18nIndexes[lang][mod]; !has {
		log.Warnf("index add failed: mod: %v, expect: %v", mod, i.modIds)
		return
	}

	i.i18nIndexes[lang][mod].rangeIndex[tp] = append(i.i18nIndexes[lang][mod].rangeIndex[tp], vo)
}

func (i *MemIndexer) AddNameIndex(vo *VO) {
	mod := vo.ModId
	tp := vo.Type
	lang := vo.Lang
	name := vo.Name

	if _, has := i.i18nIndexes[lang]; !has {
		log.Warnf("lang: %v not found", lang)
		return
	}

	if _, has := i.i18nIndexes[lang][mod]; !has {
		log.Warnf("index add failed: mod: %v, expect: %v", mod, i.modIds)
		return
	}

	if _, has := i.i18nIndexes[lang][mod].nameIndexes[tp]; !has {
		i.i18nIndexes[lang][mod].nameIndexes[tp] = make(map[string][]*VO)
	}

	i.i18nIndexes[lang][mod].nameIndexes[tp][name] = append(i.i18nIndexes[lang][mod].nameIndexes[tp][name], vo)
}

func (i *MemIndexer) AddIdIndex(vo *VO) {
	mod := vo.ModId
	tp := vo.Type
	lang := vo.Lang
	id := vo.Id

	if _, has := i.i18nIndexes[lang]; !has {
		log.Warnf("lang: %v not found", lang)
		return
	}

	if _, has := i.i18nIndexes[lang][mod]; !has {
		log.Warnf("index add failed: mod: %v, expect: %v", mod, i.modIds)
		return
	}

	if _, has := i.i18nIndexes[lang][mod].idIndex[tp]; !has {
		i.i18nIndexes[lang][mod].idIndex[tp] = make(map[string][]*VO)
	}

	i.i18nIndexes[lang][mod].idIndex[tp][id] = append(i.i18nIndexes[lang][mod].idIndex[tp][id], vo)
}
