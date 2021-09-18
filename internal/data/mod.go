package data

import (
	"zenith/internal/i18n"
	"zenith/internal/loader"
	"zenith/pkg/jsonutil"

	"github.com/leonelquinteros/gotext"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Mod struct {
	ID           string
	Name         string
	Description  string
	Path         string
	Dependencies []string
	IdMap        map[string][]*gjson.Result
	NameMap      map[string][]*gjson.Result
	TempData     map[string][]*gjson.Result
	Loaded       bool
}

func (mod *Mod) GetById(id string) []*gjson.Result {
	return mod.IdMap[id]
}

func (mod *Mod) GetByName(name string) []*gjson.Result {
	return mod.NameMap[name]
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
	if type_, has := jsonutil.GetString("type", json, ""); !has {
		log.Debugf("field type does not exist, json %s", json)
		return
	} else {
		switch type_ {
		case "MONSTER":
			processMonster(json)
		}
	}

}

func processMonster(json *gjson.Result) {
	processMonsterDiff(json)
}

func processMonsterDiff(json *gjson.Result) {
	meleeSkill, _ := jsonutil.GetInt("melee_skill", json, 0)
	meleeDice, _ := jsonutil.GetInt("melee_dice", json, 0)
	bonusCut, _ := jsonutil.GetInt("melee_cut", json, 0)
	meleeDiceSides, _ := jsonutil.GetInt("melee_dice_sides", json, 0)
	dodge, _ := jsonutil.GetInt("dodge", json, 0)
	armorBash, _ := jsonutil.GetInt("armor_bash", json, 0)
	armorCut, _ := jsonutil.GetInt("armor_cut", json, 0)
	diff, _ := jsonutil.GetInt("diff", json, 0)
	specialAttacks, _ := jsonutil.GetArray("special_attacks", json, make([]gjson.Result, 0))
	specialAttacksSize := len(specialAttacks)
	emitFields, _ := jsonutil.GetArray("emit_fields", json, make([]gjson.Result, 0))
	emitFieldsSize := len(emitFields)
	hp, _ := jsonutil.GetInt("hp", json, 0)
	speed, _ := jsonutil.GetInt("speed", json, 0)
	attackCost, _ := jsonutil.GetInt("attack_cost", json, 0)
	morale, _ := jsonutil.GetInt("morale", json, 0)
	argo, _ := jsonutil.GetInt("aggression", json, 0)
	visionDay, _ := jsonutil.GetInt("vision_day", json, 0)
	visionNight, _ := jsonutil.GetInt("vision_night", json, 0)

	difficulty := float64((meleeSkill+1)*meleeDice*(bonusCut+meleeDiceSides))*0.04 +
		float64((dodge+1)*(3+armorBash+armorCut))*0.04 + float64((diff + int64(specialAttacksSize) + 8*int64(emitFieldsSize)))

	difficulty *= (float64(hp+speed-attackCost)+float64(morale+argo)*0.1)*0.01 + float64(visionDay+2*visionNight)*0.01
	res, _ := sjson.Set(json.String(), "difficulty", difficulty)

	if difficulty < 3 {
		res, _ = sjson.Set(res, "diff_color", "light_gray")
		res, _ = sjson.Set(res, "diff_desc", "<color_light_gray>Minimal threat.")
	} else if difficulty < 10 {
		res, _ = sjson.Set(res, "diff_color", "light_gray")
		res, _ = sjson.Set(res, "diff_desc", "<color_light_gray>Mildly dangerous.</color>")
	} else if difficulty < 20 {
		res, _ = sjson.Set(res, "diff_color", "light_red")
		res, _ = sjson.Set(res, "diff_desc", "<color_light_red>Dangerous.</color>")
	} else if difficulty < 30 {
		res, _ = sjson.Set(res, "diff_color", "red")
		res, _ = sjson.Set(res, "diff_desc", "<color_red>Very dangerous.</color>")
	} else if difficulty < 50 {
		res, _ = sjson.Set(res, "diff_color", "red")
		res, _ = sjson.Set(res, "diff_desc", "<color_red>Extremely dangerous.")
	} else {
		res, _ = sjson.Set(res, "diff_color", "red")
		res, _ = sjson.Set(res, "diff_desc", "<color_red>Fatally dangerous!</color>")
	}

	*json = gjson.Parse(res)
}

func (mod *Mod) createIndex(id string, json *gjson.Result, mo *gotext.Mo) {
	name := i18n.Tran("name", json, mo)
	mod.IdMap[id] = append(mod.IdMap[id], json)
	mod.NameMap[name] = append(mod.NameMap[name], json)
}
