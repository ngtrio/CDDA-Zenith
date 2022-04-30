package view

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"zenith/internal/data"
	"zenith/internal/i18n"
	"zenith/pkg/jsonutil"

	"github.com/leonelquinteros/gotext"
	"github.com/tidwall/gjson"
)

// default value see: https://github.com/CleverRaven/Cataclysm-DDA/blob/master/src/mtype.h
func (m *VO) Bind(raw *gjson.Result, mo *gotext.Mo, po *gotext.Po) {
	m.ID, _ = jsonutil.GetString("id", raw, "")
	m.Type, _ = jsonutil.GetString("type", raw, "")
	m.Name = i18n.Tran("name", raw, mo)
	m.Description = i18n.Tran("description", raw, mo)

	m.SymbolColor, _ = jsonutil.GetString("color", raw, "")
	m.Symbol, _ = jsonutil.GetString("symbol", raw, "")

	m.DiffColor, _ = jsonutil.GetString("diff_color", raw, "")
	temp := i18n.Tran("diff_desc", raw, mo)
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

	flags, _ := jsonutil.GetArray("flags", raw, make([]gjson.Result, 0))
	for _, flag := range flags {
		// trim "PATH_" on PATH_AVOID_DANGER_x
		m.FlagsDesc = append(m.FlagsDesc, i18n.TranUI(data.Flags[strings.TrimPrefix(flag.String(), "PATH_")], po))
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
		m.Symbol, m.Name, m.ID, m.Mod,
		m.DiffDesc, m.Diff, m.Description,
		i18n.TranUI("HP", po), m.HP, i18n.TranUI("Speed", po), m.Speed,
		i18n.TranUI("Volume", po), m.Volume, i18n.TranUI("Weight", po), m.Weight,
		i18n.TranUI("Attack", po), m.Attack, i18n.TranUI("Attack cost", po), m.AttackCost,
		i18n.TranUI("Aggression", po), m.Aggression, i18n.TranUI("Morale", po), m.Morale,
		i18n.TranUI("Vision", po), m.VisionDay, i18n.TranUI("night", po), m.VisionNight,
		i18n.TranUI("Armor bash", po), m.ArmorBash, i18n.TranUI("Armor cut", po), m.ArmorCut, i18n.TranUI("Armor stab", po), m.ArmorStab,
		i18n.TranUI("Armor bullet", po), m.ArmorBullet, i18n.TranUI("Armor acid", po), m.ArmorAcid, i18n.TranUI("Armor fire", po), m.ArmorFire,
		i18n.TranUI("Bleed rate", po), m.BleedRate,
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
