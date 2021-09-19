package view

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"zenith/internal/color"
	"zenith/internal/data"
	"zenith/internal/i18n"
	"zenith/pkg/jsonutil"

	"github.com/leonelquinteros/gotext"
	"github.com/tidwall/gjson"
)

func (m *Monster) Bind(raw *gjson.Result, mo *gotext.Mo, po *gotext.Po) {
	m.ID, _ = jsonutil.GetString("id", raw, "")
	m.Type, _ = jsonutil.GetString("type", raw, "")
	m.Name = i18n.Tran("name", raw, mo)
	m.Desc = i18n.Tran("description", raw, mo)

	cr, _ := jsonutil.GetString("color", raw, "")
	symbol, _ := jsonutil.GetString("symbol", raw, "")
	var colorLoader color.Color
	colorLoader.Load(cr)
	m.Symbol = colorLoader.Colorized(symbol)

	diffColor, _ := jsonutil.GetString("diff_color", raw, "")
	temp := i18n.Tran("diff_desc", raw, mo)
	l := strings.Index(temp, ">")
	r := strings.LastIndex(temp, "<")
	diffDesc := temp[l+1 : r]
	colorLoader.Load(diffColor)
	m.DiffDesc = colorLoader.Colorized(diffDesc)

	value, _ := jsonutil.GetFloat("difficulty", raw, 0)
	m.Diff, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", value), 64)

	meleeCut, _ := jsonutil.GetInt("melee_cut", raw, 0)
	meleeDice, _ := jsonutil.GetInt("melee_dice", raw, 0)
	meleeDiceSides, _ := jsonutil.GetInt("melee_dice_sides", raw, 0)
	m.Attack = fmt.Sprintf("%dd%d+%d", meleeDice, meleeDiceSides, meleeCut)

	m.ArmorCut, _ = jsonutil.GetInt("armor_cut", raw, 0)
	m.ArmorBash, _ = jsonutil.GetInt("armor_bash", raw, 0)
	m.ArmorBullet, _ = jsonutil.GetInt("armor_bullet", raw, 0)
	m.ArmorStab, _ = jsonutil.GetInt("armor_stab", raw, 0)
	m.ArmorAcid, _ = jsonutil.GetInt("armor_acid", raw, 0)
	m.ArmorFire, _ = jsonutil.GetInt("armor_fire", raw, 0)

	m.Hp, _ = jsonutil.GetInt("hp", raw, 0)
	m.Speed, _ = jsonutil.GetInt("speed", raw, 0)

	m.Volume, _ = jsonutil.GetString("volume", raw, "0 ml")
	m.Weight, _ = jsonutil.GetString("weight", raw, "0 kg")

	m.VisionDay, _ = jsonutil.GetInt("vision_day", raw, 0)
	m.VisionNight, _ = jsonutil.GetInt("vision_night", raw, 0)

	flags, _ := jsonutil.GetArray("flags", raw, make([]gjson.Result, 0))
	for _, flag := range flags {
		m.FlagsDesc = append(m.FlagsDesc, i18n.TranUI(data.Flags[flag.String()], po))
	}
}

func (m *Monster) CliView(po *gotext.Po) string {
	template := `
%s %s(%s)
---
%s
---
%s(%.3f) 
---
%s
---
%s: %s
%s: %d	%s: %d
%s: %s	%s: %s
%s: %d	%s: %d
---
%s: %d	%s: %d	%s: %d
%s: %d	%s: %d	%s: %d
---
`
	res := fmt.Sprintf(template,
		m.Symbol, m.Name, m.ID, m.Mod,
		m.DiffDesc, m.Diff, m.Desc,
		i18n.TranUI("Attack", po), m.Attack,
		i18n.TranUI("HP", po), m.Hp, i18n.TranUI("Speed", po), m.Speed,
		i18n.TranUI("Volume", po), m.Volume, i18n.TranUI("Weight", po), m.Weight,
		i18n.TranUI("Vision day", po), m.VisionDay, i18n.TranUI("Vision night", po), m.VisionNight,
		i18n.TranUI("Armor bash", po), m.ArmorBash, i18n.TranUI("Armor cut", po), m.ArmorCut, i18n.TranUI("Armor stab", po), m.ArmorStab,
		i18n.TranUI("Armor bullet", po), m.ArmorBullet, i18n.TranUI("Armor acid", po), m.ArmorAcid, i18n.TranUI("Armor fire", po), m.ArmorFire,
	)

	for _, flagDesc := range m.FlagsDesc {
		res += "* " + flagDesc + "\n"
	}

	return res
}

func (m *Monster) JsonView() string {
	bytes, _ := json.Marshal(m)
	return string(bytes)
}
