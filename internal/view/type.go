package view

import (
	"github.com/leonelquinteros/gotext"
	"github.com/tidwall/gjson"
)

type Type interface {
	Bind(json *gjson.Result, mo *gotext.Mo, po *gotext.Po)
	CliView(po *gotext.Po) string
	JsonView() string
}

type BaseType struct {
	Mod    string `json:"mod"`
	ID     string `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Desc   string `json:"description"`
	Symbol string `json:"symbol"`
}

type Monster struct {
	BaseType
	AttackCost  int64    `json:"attack_cost"`
	BleedRate   int64    `json:"bleed_rate"`
	DiffDesc    string   `json:"diff_desc"`
	Diff        float64  `json:"difficulty"`
	Volume      string   `json:"volume"`
	Weight      string   `json:"weight"`
	Hp          int64    `json:"hp"`
	Speed       int64    `json:"speed"`
	Attack      string   `json:"attack"`
	Aggression  int64    `json:"aggression"`
	Morale      int64    `json:"morale"`
	ArmorBash   int64    `json:"armor_bash"`
	ArmorCut    int64    `json:"armor_cut"`
	ArmorBullet int64    `json:"armor_bullet"`
	ArmorStab   int64    `json:"armor_stab"`
	ArmorAcid   int64    `json:"armor_acid"`
	ArmorFire   int64    `json:"armor_fire"`
	VisionDay   int64    `json:"vision_day"`
	VisionNight int64    `json:"vision_night"`
	FlagsDesc   []string `json:"flags_description"`
}
