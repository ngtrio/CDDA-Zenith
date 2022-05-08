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
	"zenith/pkg/jsonutil"
)

type Type interface {
	Bind(tp string, json *gjson.Result, mo *gotext.Mo, po *gotext.Po)
	CliView(po *gotext.Po) string
	JsonView() string
}

type BaseType struct {
	Lang        string `json:"lang"`
	ModId       string `json:"mod_id"`
	ModName     string `json:"mod_name"`
	Id          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Symbol      string `json:"symbol"`
	SymbolColor string `json:"symbol_color"`
}

type Monster struct {
	DiffColor      string           `json:"diff_color"`
	AttackCost     int64            `json:"attack_cost"`
	BleedRate      int64            `json:"bleed_rate"`
	DiffDesc       string           `json:"diff_desc"`
	Diff           float64          `json:"difficulty"`
	Volume         string           `json:"volume"`
	Weight         string           `json:"weight"`
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
	FlagsDesc      []string         `json:"flags_description"`
	SpecialAttacks []*monsterAttack `json:"special_attacks"`
}

type MonsterAttack struct {
	Cooldown int64                  `json:"cooldown"`
	MoveCost int64                  `json:"move_cost"`
	Effects  []*monsterAttackEffect `json:"effects"`
}

type monsterAttack struct {
	Id string `json:"id"`
	MonsterAttack
}

type monsterAttackEffect struct {
	Id       string `json:"id"`
	Duration int    `json:"duration"`
	Effect
}

type Effect struct {
	Names []string `json:"name"`
	Descs []string `json:"desc"`
}

type VO struct {
	BaseType
	Monster
	MonsterAttack
	Effect
}

func NewVO(modId, modName string) *VO {
	return &VO{
		BaseType: BaseType{
			ModId:   modId,
			ModName: modName,
		},
	}
}

// default value see: https://github.com/CleverRaven/Cataclysm-DDA/blob/master/src/mtype.h
func (m *VO) Bind(raw *gjson.Result, langPack LangPack, mod *Mod, indexer Indexer) {
	tp := getType(raw)

	m.bindCommon(raw, langPack)

	// TODO add new type, step 2
	switch tp {
	case constdef.TypeMonster:
		m.bindMonster(raw, langPack, mod, indexer)
	case constdef.TypeMonsterAttack:
		m.bindMonsterAttack(raw, langPack, mod, indexer)
	case constdef.TypeEffect:
		m.bindEffect(raw, langPack, mod, indexer)
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
	colorLoader.Load(m.DiffColor)
	m.DiffDesc = colorLoader.Colorized(m.DiffDesc)

	res := fmt.Sprintf(template,
		m.Symbol, m.Name, m.Id, m.ModName,
		m.DiffDesc, m.Diff, m.Description,
		i18n.TranCustom("HP", po), m.HP, i18n.TranCustom("Speed", po), m.Speed,
		i18n.TranCustom("Volume", po), m.Volume, i18n.TranCustom("Weight", po), m.Weight,
		i18n.TranCustom("Attack", po), m.Attack, i18n.TranCustom("Attack cost", po), m.AttackCost,
		i18n.TranCustom("Aggression", po), m.Aggression, i18n.TranCustom("Morale", po), m.Morale,
		i18n.TranCustom("Vision", po), m.VisionDay, i18n.TranCustom("night", po), m.VisionNight,
		i18n.TranCustom("Armor bash", po), m.ArmorBash, i18n.TranCustom("Armor cut", po), m.ArmorCut, i18n.TranCustom("Armor stab", po), m.ArmorStab,
		i18n.TranCustom("Armor bullet", po), m.ArmorBullet, i18n.TranCustom("Armor acid", po), m.ArmorAcid, i18n.TranCustom("Armor fire", po), m.ArmorFire,
		i18n.TranCustom("Bleed rate", po), m.BleedRate,
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

func (m *VO) bindCommon(raw *gjson.Result, langPack LangPack) {
	m.Lang = langPack.Lang
	m.Id = getId(raw)
	m.Type = getType(raw)

	m.Name = i18n.Tran("name", raw, langPack.Mo)
	if m.Name == "" {
		m.Name = i18n.TranCustom(m.Id, langPack.Po)
	}

	m.Description = i18n.Tran("description", raw, langPack.Mo)
	m.ModName = i18n.TranString(m.ModName, langPack.Mo)

	m.SymbolColor, _ = jsonutil.GetString("color", raw, "")
	m.Symbol, _ = jsonutil.GetString("symbol", raw, "")
	flags, _ := jsonutil.GetArray("flags", raw, make([]gjson.Result, 0))
	for _, flag := range flags {
		// trim "PATH_" on PATH_AVOID_DANGER_x
		m.FlagsDesc = append(m.FlagsDesc, i18n.TranCustom(data.Flags[strings.TrimPrefix(flag.String(), "PATH_")], langPack.Po))
	}
}

func (m *VO) bindMonster(raw *gjson.Result, langPack LangPack, mod *Mod, indexer Indexer) {
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
	volume, _ := jsonutil.GetString("volume", raw, "0 ml")
	weight, _ := jsonutil.GetString("weight", raw, "0 kg")

	m.AttackCost = attackCost
	m.BleedRate = bleedRate
	m.Attack = fmt.Sprintf("%dd%d+%d", meleeDice, meleeDiceSides, meleeCut)
	m.Aggression = aggression
	m.Morale = morale
	m.ArmorCut = armorCut
	m.ArmorBash = armorBash
	m.ArmorBullet = armorBullet
	m.ArmorStab = armorStab
	m.ArmorAcid = armorAcid
	m.ArmorFire = armorFire
	m.HP = int64(math.Max(1, float64(hp)))
	m.Speed = speed
	m.Volume = volume
	m.Weight = weight
	m.VisionDay = visionDay
	m.VisionNight = visionNight

	// https://github.com/CleverRaven/Cataclysm-DDA/blob/c5953acae3bb4a0b2b51ddf23ea695f41079d2a8/src/monstergenerator.cpp#L1064
	m.Diff = float64((meleeSkill+1)*meleeDice*(bonusCut+meleeDiceSides))*0.04 +
		float64((dodge+1)*(3+armorBash+armorCut))*0.04 + float64(diff+int64(specialAttacksSize)+8*int64(emitFieldsSize))
	m.Diff *= (float64(hp+speed-attackCost)+float64(morale+aggression)*0.1)*0.01 + float64(visionDay+2*visionNight)*0.01
	m.Diff, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", m.Diff), 64)

	if m.Diff < 3 {
		m.DiffColor = "light_gray"
		m.DiffDesc = "<color_light_gray>Minimal threat.</color>"
	} else if m.Diff < 10 {
		m.DiffColor = "light_gray"
		m.DiffDesc = "<color_light_gray>Mildly dangerous.</color>"
	} else if m.Diff < 20 {
		m.DiffColor = "light_red"
		m.DiffDesc = "<color_light_red>Dangerous.</color>"
	} else if m.Diff < 30 {
		m.DiffColor = "red"
		m.DiffDesc = "<color_red>Very dangerous.</color>"
	} else if m.Diff < 50 {
		m.DiffColor = "red"
		m.DiffDesc = "<color_red>Extremely dangerous.</color>"
	} else {
		m.DiffColor = "red"
		m.DiffDesc = "<color_red>Fatally dangerous!</color>"
	}

	temp := i18n.TranString(m.DiffDesc, langPack.Mo)
	l := strings.Index(temp, ">")
	r := strings.LastIndex(temp, "<")
	m.DiffDesc = temp[l+1 : r]

	m.SpecialAttacks = parseSpecialAttacks(raw.Get("special_attacks"), langPack, indexer)
}

func parseSpecialAttacks(field gjson.Result, langPack LangPack, indexer Indexer) []*monsterAttack {
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
				ref := indexer.IdIndex(constdef.TypeMonsterAttack, attackId, langPack.Lang)
				ma.MonsterAttack = ref[0].MonsterAttack
			}
		}
		mas = append(mas, ma)
	}

	return mas
}

func (m *VO) bindMonsterAttack(raw *gjson.Result, langPack LangPack, mod *Mod, indexer Indexer) {
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
			es := indexer.IdIndex(constdef.TypeEffect, effectId, langPack.Lang)
			e := es[0]
			effects[i].Effect = e.Effect
		}
	}
	m.MonsterAttack = MonsterAttack{
		Cooldown: cooldown,
		MoveCost: moveCost,
		Effects:  effects,
	}
}

func (m *VO) bindEffect(raw *gjson.Result, langPack LangPack, mod *Mod, indexer Indexer) {
	e := new(Effect)
	_ = json.Unmarshal([]byte(raw.String()), e)
	names := make([]string, 0)
	descs := make([]string, 0)
	for _, name := range e.Names {
		names = append(names, i18n.TranString(name, langPack.Mo))
	}
	for _, desc := range e.Descs {
		descs = append(descs, i18n.TranString(desc, langPack.Mo))
	}

	m.Names = names
	m.Descs = descs
}
