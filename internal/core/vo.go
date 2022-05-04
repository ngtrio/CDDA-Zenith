package core

import (
	"encoding/json"
	"fmt"
	"github.com/leonelquinteros/gotext"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
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
	SpecialAttacks []*MonsterAttack `json:"special_attacks"`
}

type MonsterAttack struct {
	AttackType string                 `json:"attack_type"`
	Cooldown   int                    `json:"cooldown"`
	MoveCost   int                    `json:"move_cost"`
	Effects    []*monsterAttackEffect `json:"effects"`
}

type monsterAttackEffect struct {
	ID       string `json:"id"`
	Duration int    `json:"duration"`
	Effect   Effect `json:"effect"`
}

type Effect struct {
	Names []string `json:"names"`
	Desc  []string `json:"desc"`
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
func (m *VO) Bind(raw *gjson.Result, langPack LangPack, mod *Mod) {
	tp := getType(raw)

	m.bindCommon(raw, langPack)

	// TODO add new type, step 2
	switch tp {
	case constdef.TypeMonster:
		m.bindMonster(raw, langPack, mod)
	//case constdef.TypeMonsterAttack:
	//	m.bindMonsterAttack(raw, langPack, mod)
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

func (m *VO) bindMonster(raw *gjson.Result, langPack LangPack, mod *Mod) {
	m.DiffColor, _ = jsonutil.GetString("diff_color", raw, "")
	temp := i18n.Tran("diff_desc", raw, langPack.Mo)
	l := strings.Index(temp, ">")
	r := strings.LastIndex(temp, "<")
	diffDesc := temp[l+1 : r]

	m.DiffDesc = diffDesc

	value, _ := jsonutil.GetFloat("difficulty", raw, 0)
	m.Diff, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", value), 64)

	m.AttackCost, _ = jsonutil.GetInt("attack_cost", raw, 100)
	m.BleedRate, _ = jsonutil.GetInt("bleed_rate", raw, 100)

	meleeCut, _ := jsonutil.GetInt("melee_cut", raw, 0)
	meleeDice, _ := jsonutil.GetInt("melee_dice", raw, 0)
	meleeDiceSides, _ := jsonutil.GetInt("melee_dice_sides", raw, 0)
	m.Attack = fmt.Sprintf("%dd%d+%d", meleeDice, meleeDiceSides, meleeCut)
	m.Aggression, _ = jsonutil.GetInt("aggression", raw, 0)
	m.Morale, _ = jsonutil.GetInt("morale", raw, 0)

	m.ArmorCut, _ = jsonutil.GetInt("armor_cut", raw, -1)
	m.ArmorBash, _ = jsonutil.GetInt("armor_bash", raw, -1)
	m.ArmorBullet, _ = jsonutil.GetInt("armor_bullet", raw, -1)
	m.ArmorStab, _ = jsonutil.GetInt("armor_stab", raw, -1)
	m.ArmorAcid, _ = jsonutil.GetInt("armor_acid", raw, -1)
	m.ArmorFire, _ = jsonutil.GetInt("armor_fire", raw, -1)

	m.HP, _ = jsonutil.GetInt("hp", raw, 0)
	m.Speed, _ = jsonutil.GetInt("speed", raw, 0)

	m.Volume, _ = jsonutil.GetString("volume", raw, "0 ml")
	m.Weight, _ = jsonutil.GetString("weight", raw, "0 kg")

	m.VisionDay, _ = jsonutil.GetInt("vision_day", raw, 40)
	m.VisionNight, _ = jsonutil.GetInt("vision_night", raw, 1)

	m.SpecialAttacks = parseSpecialAttacks(raw.Get("special_attacks"))
}

func parseSpecialAttacks(field gjson.Result) []*MonsterAttack {
	return nil
}

func (m *VO) bindCommon(raw *gjson.Result, langPack LangPack) {
	m.Lang = langPack.Lang
	m.Id, _ = jsonutil.GetString("id", raw, "")
	m.Type, _ = jsonutil.GetString("type", raw, "")

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

func (m *VO) bindMonsterAttack(raw *gjson.Result, langPack LangPack, mod *Mod) {
	m.MonsterAttack = *packMonsterAttack(raw, langPack, mod)
}

func packMonsterAttack(raw *gjson.Result, langPack LangPack, mod *Mod) *MonsterAttack {
	attackType, _ := jsonutil.GetString("attack_type", raw, "")
	cooldown, _ := jsonutil.GetInt("cooldown", raw, 0)
	moveCost, _ := jsonutil.GetInt("move_cost", raw, 0)

	var effects []*monsterAttackEffect
	effectJsons, _ := jsonutil.GetArray("effects", raw, nil)
	if effectJsons != nil {
		effects = make([]*monsterAttackEffect, 0, len(effectJsons))
		for i, j := range effectJsons {
			_ = json.Unmarshal([]byte(j.String()), effects[i])
			//mod.GetById()
		}
	}

	return &MonsterAttack{
		AttackType: attackType,
		Cooldown:   int(cooldown),
		MoveCost:   int(moveCost),
		Effects:    effects,
	}
}
