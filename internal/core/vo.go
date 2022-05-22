package core

import (
	"encoding/json"
	"fmt"
	"github.com/leonelquinteros/gotext"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"math"
	"strconv"
	"strings"
	"zenith/internal/constdef"
	"zenith/internal/data"
	"zenith/internal/i18n"
	"zenith/internal/util"
	"zenith/pkg/jsonutil"
)

type Type interface {
	Bind(tp string, json *gjson.Result, mo *gotext.Mo, po *gotext.Po)
	CliView(po *gotext.Po) string
	JsonView() string
}

type baseType struct {
	Lang        string   `json:"lang"`
	ModId       string   `json:"mod_id"`
	ModName     string   `json:"mod_name"`
	Id          string   `json:"id"`
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Symbol      string   `json:"symbol"`
	SymbolColor string   `json:"symbol_color"`
	Volume      string   `json:"volume"`
	Weight      string   `json:"weight"`
	Material    []string `json:"material"`
	FlagsDesc   []string `json:"flags_description"`
}

type Monster struct {
	DiffColor      string           `json:"diff_color"`
	AttackCost     int64            `json:"attack_cost"`
	BleedRate      int64            `json:"bleed_rate"`
	DiffDesc       string           `json:"diff_desc"`
	Diff           float64          `json:"difficulty"`
	HP             int64            `json:"hp"`
	Speed          int64            `json:"speed"`
	Attack         string           `json:"attack"`
	Aggression     int64            `json:"aggression"`
	Morale         int64            `json:"morale"`
	ArmorBash      int64            `json:"armor_bash"`
	ArmorCut       int64            `json:"armor_cut"`
	ArmorBullet    int64            `json:"armor_bullet"`
	ArmorStab      int64            `json:"armor_stab"`
	ArmorAcid      int64            `json:"armor_acid"`
	ArmorFire      int64            `json:"armor_fire"`
	VisionDay      int64            `json:"vision_day"`
	VisionNight    int64            `json:"vision_night"`
	SpecialAttacks []*monsterAttack `json:"special_attacks"`
}

type MonsterAttack struct {
	Cooldown int64                  `json:"cooldown"`
	MoveCost int64                  `json:"move_cost"`
	Effects  []*monsterAttackEffect `json:"effects"`
}

type monsterAttack struct {
	Id    string `json:"id"`
	ModId string `json:"mod_id"`
	MonsterAttack
}

type monsterAttackEffect struct {
	ModId    string `json:"mod_id"`
	Id       string `json:"id"`
	Duration int    `json:"duration"`
	Effect
}

type Effect struct {
	Names []string `json:"name"`
	Descs []string `json:"desc"`
}

type Item struct {
	baseType
	baseItem
	Replace string `json:"replace"`
}

type baseItem struct {
	Price     int64             `json:"price"`
	Qualities []quality         `json:"qualities"`
	CraftFrom map[string][]item `json:"craft_from"`
	UnCraftTo map[string][]item `json:"un_craft_to"`
}

type sub struct {
	ModId string `json:"modId"`
	Id    string `json:"id"`
	Name  string `json:"name"`
	SType string `json:"s_type"`
}

type item struct {
	sub
	Num int64 `json:"num"`
}

type quality struct {
	sub
	Level int `json:"level"`
}

type Requirement struct {
	Components [][]item  `json:"components"`
	Tools      [][]item  `json:"tools"`
	Qualities  []quality `json:"qualities"`
}

type Recipe struct {
	SkillUsed      string    `json:"skill_used"`
	Difficulty     string    `json:"difficulty"`
	SkillsRequired string    `json:"skills_required"`
	Components     [][]item  `json:"components"`
	Tools          [][]item  `json:"tools"`
	Qualities      []quality `json:"qualities"`
	Using          []item    `json:"using"`
}

type VO struct {
	baseType
	Monster       *Monster       `json:"monster,omitempty"`
	MonsterAttack *MonsterAttack `json:"monster_attack,omitempty"`
	Effect        *Effect        `json:"effect,omitempty"`
	Item          *Item          `json:"item,omitempty"`
	Recipe        *Recipe        `json:"recipe,omitempty"`
	Requirement   *Requirement   `json:"requirement,omitempty"`
}

func NewVO(modId, modName string) *VO {
	return &VO{
		baseType: baseType{
			ModId:   modId,
			ModName: modName,
		},
	}
}

// default value see: https://github.com/CleverRaven/Cataclysm-DDA/blob/master/src/mtype.h
func (m *VO) Bind(raw *gjson.Result, langPack LangPack, mod *Mod, game *Game) {
	m.bindCommon(raw, langPack, mod, game)

	// TODO add new type, step 2
	switch tp := getType(raw); {
	case tp == constdef.TypeMonster:
		m.bindMonster(raw, langPack, mod, game)
	case tp == constdef.TypeMonsterAttack:
		m.bindMonsterAttack(raw, langPack, mod, game)
	case tp == constdef.TypeEffect:
		m.bindEffect(raw, langPack)
	case tp == constdef.TypeRecipe || tp == constdef.TypeUnCraft:
		m.preBindRecipeAndUnCraft(raw, langPack, mod, game)
	case tp == constdef.TypeRequirement:
		m.bindRequirement(raw, langPack, mod, game)
	case constdef.ItemTypes[tp]:
		m.bindItem(raw, langPack, mod, game)
	default:
		log.Debugf("type %v is not supported yet", tp)
	}
}

func (m *VO) CliView(po *gotext.Po) string {
	template := `
%s %s(%s)
---
%s
---
%s(%.3f) 
---
%s
---
%s: %d	%s: %d
%s: %s	%s: %s
---
%s: %s	%s: %d
%s: %d	%s: %d
%s: %d (%s: %d)
---
%s: %d	%s: %d	%s: %d
%s: %d	%s: %d	%s: %d
%s: %d
---
`

	var colorLoader Color
	colorLoader.Load(m.SymbolColor)
	m.Symbol = colorLoader.Colorized(m.Symbol)
	colorLoader.Load(m.Monster.DiffColor)
	m.Monster.DiffDesc = colorLoader.Colorized(m.Monster.DiffDesc)

	res := fmt.Sprintf(template,
		m.Symbol, m.Name, m.Id, m.ModName,
		m.Monster.DiffDesc, m.Monster.Diff, m.Description,
		i18n.TranCustom("HP", po), m.Monster.HP, i18n.TranCustom("Speed", po), m.Monster.Speed,
		i18n.TranCustom("Volume", po), m.Volume, i18n.TranCustom("Weight", po), m.Weight,
		i18n.TranCustom("Attack", po), m.Monster.Attack, i18n.TranCustom("Attack cost", po), m.Monster.AttackCost,
		i18n.TranCustom("Aggression", po), m.Monster.Aggression, i18n.TranCustom("Morale", po), m.Monster.Morale,
		i18n.TranCustom("Vision", po), m.Monster.VisionDay, i18n.TranCustom("night", po), m.Monster.VisionNight,
		i18n.TranCustom("Armor bash", po), m.Monster.ArmorBash, i18n.TranCustom("Armor cut", po), m.Monster.ArmorCut, i18n.TranCustom("Armor stab", po), m.Monster.ArmorStab,
		i18n.TranCustom("Armor bullet", po), m.Monster.ArmorBullet, i18n.TranCustom("Armor acid", po), m.Monster.ArmorAcid, i18n.TranCustom("Armor fire", po), m.Monster.ArmorFire,
		i18n.TranCustom("Bleed rate", po), m.Monster.BleedRate,
	)

	for _, flagDesc := range m.FlagsDesc {
		res += "* " + flagDesc + "\n"
	}

	return res
}

func (m *VO) JsonView() string {
	bytes, _ := json.Marshal(m)
	return string(bytes)
}

func (m *VO) bindCommon(raw *gjson.Result, langPack LangPack, mod *Mod, game *Game) {
	m.Lang = langPack.Lang
	m.Id = getId(raw)
	m.Type = getType(raw)

	m.Name = i18n.Tran("name", raw, langPack.Mo)
	if m.Name == "" {
		m.Name = i18n.TranString(m.Id, langPack.Mo)
	}

	m.Description = i18n.Tran("description", raw, langPack.Mo)
	m.ModName = i18n.TranString(m.ModName, langPack.Mo)

	m.SymbolColor, _ = jsonutil.GetString("color", raw, "")
	m.Symbol, _ = jsonutil.GetString("symbol", raw, "")
	flags, _ := jsonutil.GetArray("flags", raw, make([]gjson.Result, 0))
	for _, flag := range flags {
		// trim "PATH_" on PATH_AVOID_DANGER_x
		flagDesc := fmt.Sprintf("%s(%s)", i18n.TranCustom(data.Flags[strings.TrimPrefix(flag.String(), "PATH_")], langPack.Po), flag.String())
		m.FlagsDesc = append(m.FlagsDesc, flagDesc)
	}

	materials, _ := jsonutil.GetArray("material", raw, nil)
	for _, material := range materials {
		var mId string
		if material.IsObject() {
			mId = getType(&material)
		} else {
			mId = material.String()
		}
		ms := m.getDepFromIndex(game, mod, constdef.TypeMaterial, mId, langPack)
		if len(ms) > 0 {
			m.Material = append(m.Material, ms[0].Name)
		} else {
			m.Material = append(m.Material, mId)
		}
	}

	volume, _ := jsonutil.GetString("volume", raw, "0 ml")
	weight, _ := jsonutil.GetString("weight", raw, "0 kg")
	m.Volume = volume
	m.Weight = weight
}

func (m *VO) bindMonster(raw *gjson.Result, langPack LangPack, mod *Mod, game *Game) {
	m.Monster = new(Monster)

	meleeCut, _ := jsonutil.GetInt("melee_cut", raw, 0)
	meleeDice, _ := jsonutil.GetInt("melee_dice", raw, 0)
	meleeDiceSides, _ := jsonutil.GetInt("melee_dice_sides", raw, 0)
	meleeSkill, _ := jsonutil.GetInt("melee_skill", raw, 0)
	bonusCut, _ := jsonutil.GetInt("melee_cut", raw, 0)
	dodge, _ := jsonutil.GetInt("dodge", raw, 0)
	armorBash, _ := jsonutil.GetInt("armor_bash", raw, 0)
	armorCut, _ := jsonutil.GetInt("armor_cut", raw, 0)
	armorBullet, _ := jsonutil.GetInt("armor_bullet", raw, 0)
	armorStab, _ := jsonutil.GetInt("armor_stab", raw, 0)
	armorAcid, _ := jsonutil.GetInt("armor_acid", raw, 0)
	armorFire, _ := jsonutil.GetInt("armor_fire", raw, 0)
	diff, _ := jsonutil.GetInt("diff", raw, 0)
	specialAttacks, _ := jsonutil.GetArray("special_attacks", raw, make([]gjson.Result, 0))
	specialAttacksSize := len(specialAttacks)
	emitFields, _ := jsonutil.GetArray("emit_fields", raw, make([]gjson.Result, 0))
	emitFieldsSize := len(emitFields)
	hp, _ := jsonutil.GetInt("hp", raw, 0)
	speed, _ := jsonutil.GetInt("speed", raw, 0)
	attackCost, _ := jsonutil.GetInt("attack_cost", raw, 100)
	morale, _ := jsonutil.GetInt("morale", raw, 0)
	aggression, _ := jsonutil.GetInt("aggression", raw, 0)
	visionDay, _ := jsonutil.GetInt("vision_day", raw, 40)
	visionNight, _ := jsonutil.GetInt("vision_night", raw, 1)
	bleedRate, _ := jsonutil.GetInt("bleed_rate", raw, 100)

	m.Monster.AttackCost = attackCost
	m.Monster.BleedRate = bleedRate
	m.Monster.Attack = fmt.Sprintf("%dd%d+%d", meleeDice, meleeDiceSides, meleeCut)
	m.Monster.Aggression = aggression
	m.Monster.Morale = morale
	m.Monster.ArmorCut = armorCut
	m.Monster.ArmorBash = armorBash
	m.Monster.ArmorBullet = armorBullet
	m.Monster.ArmorStab = armorStab
	m.Monster.ArmorAcid = armorAcid
	m.Monster.ArmorFire = armorFire
	m.Monster.HP = int64(math.Max(1, float64(hp)))
	m.Monster.Speed = speed
	m.Monster.VisionDay = visionDay
	m.Monster.VisionNight = visionNight

	// https://github.com/CleverRaven/Cataclysm-DDA/blob/c5953acae3bb4a0b2b51ddf23ea695f41079d2a8/src/monstergenerator.cpp#L1064
	m.Monster.Diff = float64((meleeSkill+1)*meleeDice*(bonusCut+meleeDiceSides))*0.04 +
		float64((dodge+1)*(3+armorBash+armorCut))*0.04 + float64(diff+int64(specialAttacksSize)+8*int64(emitFieldsSize))
	m.Monster.Diff *= (float64(hp+speed-attackCost)+float64(morale+aggression)*0.1)*0.01 + float64(visionDay+2*visionNight)*0.01
	m.Monster.Diff, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", m.Monster.Diff), 64)

	if m.Monster.Diff < 3 {
		m.Monster.DiffColor = "light_gray"
		m.Monster.DiffDesc = "<color_light_gray>Minimal threat.</color>"
	} else if m.Monster.Diff < 10 {
		m.Monster.DiffColor = "light_gray"
		m.Monster.DiffDesc = "<color_light_gray>Mildly dangerous.</color>"
	} else if m.Monster.Diff < 20 {
		m.Monster.DiffColor = "light_red"
		m.Monster.DiffDesc = "<color_light_red>Dangerous.</color>"
	} else if m.Monster.Diff < 30 {
		m.Monster.DiffColor = "red"
		m.Monster.DiffDesc = "<color_red>Very dangerous.</color>"
	} else if m.Monster.Diff < 50 {
		m.Monster.DiffColor = "red"
		m.Monster.DiffDesc = "<color_red>Extremely dangerous.</color>"
	} else {
		m.Monster.DiffColor = "red"
		m.Monster.DiffDesc = "<color_red>Fatally dangerous!</color>"
	}

	temp := i18n.TranString(m.Monster.DiffDesc, langPack.Mo)
	l := strings.Index(temp, ">")
	r := strings.LastIndex(temp, "<")
	m.Monster.DiffDesc = temp[l+1 : r]

	m.Monster.SpecialAttacks = m.parseSpecialAttacks(raw.Get("special_attacks"), langPack, mod, game)
}

func (m *VO) parseSpecialAttacks(field gjson.Result, langPack LangPack, mod *Mod, game *Game) []*monsterAttack {
	var mas []*monsterAttack
	for _, sa := range field.Array() {
		ma := &monsterAttack{}
		if sa.IsArray() {
			ma.Id = sa.Array()[0].String()
			ma.Cooldown = sa.Array()[1].Int()
		}

		if sa.IsObject() {
			if val, has := jsonutil.GetString("type", &sa, ""); has {
				ma.Id = val
				ma.Cooldown, _ = jsonutil.GetInt("cooldown", &sa, 0)
			} else {
				attackId, _ := jsonutil.GetString("id", &sa, "")
				ma.Id = attackId

				ref := m.getDepFromIndex(game, mod, constdef.TypeMonsterAttack, attackId, langPack)
				if len(ref) == 0 {
					continue
				}
				ma.MonsterAttack = *ref[0].MonsterAttack
				ma.ModId = ref[0].ModId
			}
		}
		mas = append(mas, ma)
	}

	return mas
}

func (m *VO) bindMonsterAttack(raw *gjson.Result, langPack LangPack, mod *Mod, game *Game) {
	cooldown, _ := jsonutil.GetInt("cooldown", raw, 0)
	moveCost, _ := jsonutil.GetInt("move_cost", raw, 0)

	var effects []*monsterAttackEffect
	effectJsons, _ := jsonutil.GetArray("effects", raw, nil)
	if effectJsons != nil {
		effects = make([]*monsterAttackEffect, len(effectJsons))
		for i, j := range effectJsons {
			effects[i] = new(monsterAttackEffect)
			_ = json.Unmarshal([]byte(j.String()), effects[i])
			effectId := effects[i].Id
			es := m.getDepFromIndex(game, mod, constdef.TypeEffect, effectId, langPack)
			if len(es) == 0 {
				continue
			}
			e := es[0]
			effects[i].Effect = *e.Effect
			effects[i].ModId = e.ModId
		}
	}
	m.MonsterAttack = &MonsterAttack{
		Cooldown: cooldown,
		MoveCost: moveCost,
		Effects:  effects,
	}
}

func (m *VO) bindEffect(raw *gjson.Result, langPack LangPack) {
	m.Effect = new(Effect)
	e := new(Effect)
	_ = json.Unmarshal([]byte(raw.String()), e)
	e.Names = util.Set(e.Names)
	e.Descs = util.Set(e.Descs)

	names := make([]string, 0)
	descs := make([]string, 0)
	for _, name := range e.Names {
		names = append(names, i18n.TranString(name, langPack.Mo))
	}
	for _, desc := range e.Descs {
		descs = append(descs, i18n.TranString(desc, langPack.Mo))
	}

	m.Name = i18n.TranCustom(m.Id, langPack.Po)
	m.Effect.Names = names
	m.Effect.Descs = descs
}

func (m *VO) bindRequirement(raw *gjson.Result, langPack LangPack, mod *Mod, game *Game) {
	m.Requirement = new(Requirement)

	// qualities
	m.Requirement.Qualities = m.loadQualities(raw, langPack, mod, game)

	// tools
	tools := m.preLoadItems(raw, "tools")

	// components
	components := m.preLoadItems(raw, "components")

	m.Requirement.Tools, m.Requirement.Components = tools, components
	return
}

func (m *VO) postBindRequirement(game *Game, langPack LangPack) {
	tools := m.Requirement.Tools
	components := m.Requirement.Components

	m.Requirement.Tools = m.postLoadItems(tools, langPack, "tools", game)
	m.Requirement.Components = m.postLoadItems(components, langPack, "component", game)
}

func (m *VO) preBindRecipeAndUnCraft(raw *gjson.Result, langPack LangPack, mod *Mod, game *Game) {
	m.Recipe = new(Recipe)

	m.Recipe.Tools = m.preLoadItems(raw, "tools")
	m.Recipe.Components = m.preLoadItems(raw, "components")
	m.Recipe.Qualities = m.loadQualities(raw, langPack, mod, game)
	m.Recipe.Using = m.preLoadUsing(raw)

	return
}

func (m *VO) postBindRecipeAndUnCraft(game *Game, langPack LangPack) {

	usingTools := m.postLoadUsing(m.Recipe.Using, langPack, "tools", game)
	usingComponents := m.postLoadUsing(m.Recipe.Using, langPack, "components", game)

	tools := m.postLoadItems(m.Recipe.Tools, langPack, "tools", game)
	components := m.postLoadItems(m.Recipe.Components, langPack, "components", game)

	m.Recipe.Tools = append(tools, usingTools...)
	m.Recipe.Components = append(components, usingComponents...)

	for idx, tool := range m.Recipe.Tools {
		m.Recipe.Tools[idx] = util.Set(tool)
		var newTools []item
		for _, t := range m.Recipe.Tools[idx] {
			newTools = append(newTools, t)
			toolId := t.Id
			subTools := game.ToolSub[toolId]
			for _, stId := range subTools {
				for _, st := range game.Indexer.IdIndex(constdef.TypeTool, stId, langPack.Lang) {
					it := item{
						sub: sub{
							ModId: st.ModId,
							Id:    st.Id,
							Name:  st.Name,
							SType: st.Type,
						},
						Num: t.Num,
					}
					newTools = append(newTools, it)
				}
			}
		}
		m.Recipe.Tools[idx] = newTools
	}

	for idx, component := range m.Recipe.Components {
		m.Recipe.Components[idx] = util.Set(component)
	}
}

func (m *VO) bindItem(raw *gjson.Result, langPack LangPack, mod *Mod, game *Game) {
	m.Item = new(Item)
	m.Item.CraftFrom = make(map[string][]item)
	m.Item.UnCraftTo = make(map[string][]item)

	if m.Type == constdef.TypeMigration {
		m.Item.Replace, _ = jsonutil.GetString("replace", raw, "")
		return
	}

	if m.Type == constdef.TypeTool {
		subTool, _ := jsonutil.GetString("sub", raw, "")
		if subTool != "" {
			game.ToolSub[subTool] = append(game.ToolSub[subTool], m.Id)
		}
	}

	// price
	m.Item.Price, _ = jsonutil.GetInt("price", raw, 0)

	// qualities
	m.Item.Qualities = m.loadQualities(raw, langPack, mod, game)

	// craft_from
	recipes := game.Indexer.IdIndex(constdef.TypeRecipe, m.Id, langPack.Lang)
	for _, recipe := range recipes {
		modName := recipe.ModName
		m.Item.CraftFrom[modName] = append(m.Item.CraftFrom[modName], item{
			sub: sub{
				ModId: modName,
				Id:    recipe.Id,
				Name:  recipe.Name,
				SType: recipe.Type,
			},
		})
	}

	// un_craft_to
	uncrafts := game.Indexer.IdIndex(constdef.TypeUnCraft, m.Id, langPack.Lang)
	for _, uncraft := range uncrafts {
		modName := uncraft.ModName
		m.Item.CraftFrom[modName] = append(m.Item.UnCraftTo[modName], item{
			sub: sub{
				ModId: modName,
				Id:    uncraft.Id,
				Name:  uncraft.Name,
				SType: uncraft.Type,
			},
		})
	}
}

func (m *VO) loadQualities(raw *gjson.Result, langPack LangPack, mod *Mod, game *Game) []quality {
	ql, _ := jsonutil.GetArray("qualities", raw, nil)

	var (
		id    string
		level int64
		res   []quality
	)

	for _, q := range ql {
		if q.IsObject() {
			id = getId(&q)
			level, _ = jsonutil.GetInt("level", &q, 0)
		} else if q.IsArray() {
			if q.Array()[0].IsObject() {
				id, _ = jsonutil.GetString("id", &q.Array()[0], "")
				level, _ = jsonutil.GetInt("level", &q.Array()[0], 0)
			} else {
				id = q.Array()[0].String()
				level = q.Array()[1].Int()
			}
		} else {
			log.Debugf("quality format invalid, %v, type: %v, json: %v", q, q.Type, raw)
		}

		vos := m.getDepFromIndex(game, mod, constdef.TypeToolQuality, id, langPack)
		if len(vos) == 0 {
			res = append(res, quality{
				sub: sub{
					Id:   id,
					Name: id,
				},
				Level: int(level),
			})
		} else {
			res = append(res, quality{
				sub: sub{
					ModId: vos[0].ModId,
					Id:    vos[0].Id,
					Name:  vos[0].Name,
					SType: vos[0].Type,
				},
				Level: int(level),
			})
		}
	}

	return res
}

func (m *VO) preLoadItems(raw *gjson.Result, field string) [][]item {
	var items [][]item
	itemGroups, _ := jsonutil.GetArray(field, raw, []gjson.Result{})
	for _, it := range itemGroups {
		var itemAlt []item
		for _, t := range it.Array() {
			id := t.Array()[0].String()
			num := int64(-1)
			if len(t.Array()) > 1 {
				num = t.Array()[1].Int()
			}
			var sType string
			if len(t.Array()) == 3 && t.Array()[2].String() == "LIST" {
				sType = constdef.TypeRequirement
			}

			itemAlt = append(itemAlt, item{
				sub: sub{
					Id:    id,
					SType: sType,
				},
				Num: num,
			})
		}
		items = append(items, itemAlt)
	}

	return items
}

func (m *VO) preLoadUsing(raw *gjson.Result) []item {
	var items []item

	usings, _ := jsonutil.GetArray("using", raw, []gjson.Result{})
	for _, using := range usings {
		id := using.Array()[0].String()
		num := int64(-1)
		if len(using.Array()) > 1 {
			num = using.Array()[1].Int()
		}

		items = append(items, item{
			sub: sub{
				Id:    id,
				SType: constdef.TypeRequirement,
			},
			Num: num,
		})
	}

	return items
}

func (m *VO) postLoadItems(items [][]item, langPack LangPack, field string, game *Game) [][]item {
	var newItems [][]item

	for _, it := range items {
		var itemAlt []item
		for _, t := range it {
			id := t.Id
			num := t.Num

			if t.SType == constdef.TypeRequirement {
				reqs := game.Indexer.IdIndex(t.SType, id, langPack.Lang)
				toAdd := make([]item, 0)
				for _, req := range reqs {
					var reqItem [][]item
					switch field {
					case "tools":
						reqItem = req.Requirement.Tools
					case "components":
						reqItem = req.Requirement.Components
					}
					if len(reqItem) > 1 {
						log.Warnf("loadItems LIST > 1, id: %v", id)
						return nil
					}
					if len(reqItem) == 1 {
						for _, i := range reqItem[0] {
							if i.Num != -1 {
								i.Num *= num
							}
							toAdd = append(toAdd, i)
						}
					}
				}
				itemAlt = append(itemAlt, toAdd...)

			} else {
				itemAlt = append(itemAlt, m.postLoadItem(id, num, langPack, game))
			}
		}

		newItems = append(newItems, itemAlt)
	}

	return newItems
}

func (m *VO) postLoadItem(itemId string, itemNum int64, langPack LangPack, game *Game) item {
	for tp := range constdef.ItemTypes {
		res := game.Indexer.IdIndex(tp, itemId, langPack.Lang)
		if len(res) == 0 {
			continue
		} else if tp == constdef.TypeMigration {
			replace := res[0].Item.Replace
			for tp := range constdef.ItemTypes {
				res := game.Indexer.IdIndex(tp, replace, langPack.Lang)
				if len(res) > 0 {
					break
				}
			}

			if len(res) == 0 {
				break
			}
		}

		t := res[0]
		return item{
			sub: sub{
				ModId: t.ModId,
				Id:    t.Id,
				Name:  t.Name,
				SType: t.Type,
			},
			Num: itemNum,
		}
	}

	log.Errorf("item not found, itemId: %v, id: %v", itemId, m.Id)
	return item{
		sub: sub{
			Id: itemId,
		},
		Num: itemNum,
	}
}

func (m *VO) postLoadUsing(usings []item, langPack LangPack, field string, game *Game) [][]item {
	var newItems [][]item

	for _, using := range usings {
		id := using.Id
		num := using.Num

		reqs := game.Indexer.IdIndex(constdef.TypeRequirement, id, langPack.Lang)
		for _, req := range reqs {
			var reqItems [][]item
			switch field {
			case "tools":
				reqItems = req.Requirement.Tools
			case "components":
				reqItems = req.Requirement.Components
			}

			for _, reqItem := range reqItems {
				var alt []item
				for _, i := range reqItem {
					i.Num *= num
					alt = append(alt, i)
				}
				newItems = append(newItems, alt)
			}
		}
	}

	return newItems
}

func (m *VO) getDepFromIndex(game *Game, mod *Mod, tp string, id string, langPack LangPack) []*VO {
	res := game.loadVO(mod, tp, id, langPack)
	if len(res) == 0 {
		for _, depModId := range mod.Dependencies {
			res = game.loadVO(game.Mods[depModId], tp, id, langPack)
			if len(res) != 0 {
				break
			}
		}
	}

	if len(res) == 0 {
		log.Errorf("req not found, reqId: %v, id: %v", id, m.Id)
	}

	return res
}
