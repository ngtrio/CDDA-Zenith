package constdef

import "zenith/internal/util"

// TODO add new type, step 1

var ItemTypes = map[string]bool{
	TypeAmmo:       true,
	TypeMagazine:   true,
	TypeGeneric:    true,
	TypeBattery:    true,
	TypeArmor:      true,
	TypeGun:        true,
	TypeBook:       true,
	TypeGunMod:     true,
	TypeTool:       true,
	TypeItemGroup:  true,
	TypeComestible: true,
	TypeToolArmor:  true,
}

var CommonTypes = map[string]bool{
	TypeMonster:       true,
	TypeMonsterAttack: true,
	TypeEffect:        true,
	TypeMaterial:      true,
	TypeToolQuality:   true,
	TypeRecipe:        true,
	TypeUnCraft:       true,
	TypeRequirement:   true,
}

var AllowTypeList = util.MergeMap(ItemTypes, CommonTypes)

const (
	TypeMonster       = "MONSTER"
	TypeMonsterAttack = "monster_attack"
	TypeEffect        = "effect_type"
	TypeAmmo          = "AMMO"
	TypeMagazine      = "MAGAZINE"
	TypeGeneric       = "GENERIC"
	TypeBattery       = "BATTERY"
	TypeArmor         = "ARMOR"
	TypeGun           = "GUN"
	TypeBook          = "BOOK"
	TypeGunMod        = "GUNMOD"
	TypeTool          = "TOOL"
	TypeItemGroup     = "item_group"
	TypeComestible    = "COMESTIBLE"
	TypeToolArmor     = "TOOL_ARMOR"
	TypeMaterial      = "material"
	TypeToolQuality   = "tool_quality"
	TypeRecipe        = "recipe"
	TypeUnCraft       = "uncraft"
	TypeRequirement   = "requirement"
)
