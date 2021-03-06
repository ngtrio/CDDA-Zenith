package data

var Flags = map[string]string{"SEES": "It can see you (and will run/follow).",
	"HEARS":                   "It can hear you.",
	"GOODHEARING":             "Pursues sounds more than most monsters.",
	"SMELLS":                  "It can smell you.",
	"KEENNOSE":                "Keen sense of smell.",
	"STUMBLES":                "Stumbles in its movement.",
	"WARM":                    "Warm blooded.",
	"NEMESIS":                 "Is a nemesis monster",
	"NOHEAD":                  "Headshots not allowed!",
	"HARDTOSHOOT":             "It's one size smaller for ranged attacks no less then MS_TINY",
	"GRABS":                   "Its attacks may grab you!",
	"BASHES":                  "Bashes down doors.",
	"DESTROYS":                "Bashes down walls and more. (2.5x bash multiplier where base is the critter's max melee bashing)",
	"BORES":                   "Tunnels through just about anything (15x bash multiplier: dark wyrms' bash skill 12->180)",
	"POISON":                  "Poisonous to eat.",
	"VENOM":                   "Attack may poison the player.",
	"BADVENOM":                "Attack may severely poison the player.",
	"PARALYZE":                "Attack may paralyze the player with venom.",
	"WEBWALK":                 "Doesn't destroy webs and won't get caught in them.",
	"DIGS":                    "Digs through the ground.",
	"CAN_DIG":                 "Can dig and walk.",
	"FLIES":                   "Can fly (over water etc.)",
	"AQUATIC":                 "Confined to water.",
	"SWIMS":                   "Treats water as 50 movement point terrain.",
	"ATTACKMON":               "Attacks other monsters.",
	"ANIMAL":                  "Is an animal for purposes of the Animal Empathy trait.",
	"PLASTIC":                 "Absorbs physical damage to a great degree.",
	"SUNDEATH":                "Dies in full sunlight.",
	"ELECTRIC":                "Shocks unarmed attackers.",
	"ACIDPROOF":               "Immune to acid.",
	"ACIDTRAIL":               "Leaves a trail of acid.",
	"SHORTACIDTRAIL":          "Leaves an intermittent trail of acid",
	"FIREPROOF":               "Immune to fire.",
	"SLUDGEPROOF":             "Ignores the effect of sludge trails.",
	"SLUDGETRAIL":             "Causes the monster to leave a sludge trap trail when moving.",
	"SMALLSLUDGETRAIL":        "Causes the monster to occasionally leave a 1-tile sludge trail when moving.",
	"COLDPROOF":               "Immune to cold damage.",
	"FIREY":                   "Burns stuff and is immune to fire.",
	"QUEEN":                   "When it dies local populations start to die off too.",
	"ELECTRONIC":              "e.g. A Robot; affected by emp blasts and other stuff.",
	"FUR":                     "May produce fur when butchered.",
	"LEATHER":                 "May produce leather when butchered.",
	"WOOL":                    "May produce wool when butchered.",
	"BONES":                   "May produce bones and sinews when butchered.",
	"FAT":                     "May produce fat when butchered.",
	"CONSOLE_DESPAWN":         "Despawns when a nearby console is properly hacked.",
	"IMMOBILE":                "Doesn't move (e.g. turrets)",
	"ID_CARD_DESPAWN":         "Despawns when a science ID card is used on a nearby console",
	"RIDEABLE_MECH":           "This monster is a mech suit that can be piloted.",
	"MILITARY_MECH":           "Is a military-grade mech.",
	"MECH_RECON_VISION":       "This mech grants you night-vision and enhanced overmap sight radius when piloted.",
	"MECH_DEFENSIVE":          "This mech can protect you thoroughly when piloted.",
	"HIT_AND_RUN":             "Flee for several turns after a melee attack.",
	"GUILT":                   "You feel guilty for killing it.",
	"PAY_BOT":                 "Creature can be turned into a pet for a limited time in exchange of e-money.",
	"HUMAN":                   "It's a live human as long as it's alive.",
	"NO_BREATHE":              "Creature can't drown and is unharmed by gas smoke or poison.",
	"FLAMMABLE":               "Monster catches fire burns and spreads fire to nearby objects.",
	"REVIVES":                 "Monster corpse will revive after a short period of time.",
	"CHITIN":                  "May produce chitin when butchered.",
	"VERMIN":                  "Obsolete flag for inconsequential monsters now prevents loading.",
	"NOGIB":                   "Does not leave gibs / meat chunks when killed with huge damage.",
	"LARVA":                   "Creature is a larva. Currently used for gib and blood handling.",
	"ARTHROPOD_BLOOD":         "Forces monster to bleed hemolymph.",
	"ACID_BLOOD":              "Makes monster bleed acid. Fun stuff! Does not automatically dissolve in a pool of acid on death.",
	"BILE_BLOOD":              "Makes monster bleed bile.",
	"ABSORBS":                 "Consumes objects it moves over. (Modders use this).",
	"ABSORBS_SPLITS":          "Consumes objects it moves over and if it absorbs enough it will split into a copy.",
	"CBM_CIV":                 "May produce a common and a power CBM when butchered.",
	"CBM_POWER":               "May produce a power CBM when butchered independent of other CBMs.",
	"CBM_SCI":                 "May produce a CBM from 'bionics_sci' item group when butchered.",
	"CBM_OP":                  "May produce a CBM or two from 'bionics_op' item group when butchered.",
	"CBM_TECH":                "May produce a CBM or two from 'bionics_tech' item group and a power CBM when butchered.",
	"CBM_SUBS":                "May produce a CBM or two from bionics_subs and a power CBM when butchered.",
	"FILTHY":                  "Any clothing it drops will be filthy. The squeamish trait prevents wearing clothing with this flag one can't craft anything from filthy components and wearing filthy clothes may result in infection if hit in melee.",
	"FISHABLE":                "It is fishable.",
	"GROUP_BASH":              "Gets help from monsters around it when bashing.",
	"SWARMS":                  "Groups together and forms loose packs.",
	"GROUP_MORALE":            "More courageous when near friends.",
	"INTERIOR_AMMO":           "Monster contains ammo inside itself no need to load on launch. Prevents ammo from being dropped on disable.",
	"CLIMBS":                  "Can climb over fences or similar obstacles quickly.",
	"PACIFIST":                "Monster will never do melee attacks.",
	"KEEP_DISTANCE":           "Monster will try to keep tracking_distance number of tiles between it and its current target.",
	"PUSH_MON":                "Can push creatures out of its way.",
	"PUSH_VEH":                "Monsters that can push vehicles out of their way",
	"NIGHT_INVISIBILITY":      "Monster becomes invisible if it's more than one tile away and the lighting on its tile is LL_LOW or less. Visibility is not affected by night vision.",
	"REVIVES_HEALTHY":         "When revived this monster has full hitpoints and speed",
	"NO_NECRO":                "This monster can't be revived by necros. It will still rise on its own.",
	"AVOID_DANGER_1":          "This monster will path around some dangers instead of through them.",
	"AVOID_DANGER_2":          "This monster will path around most dangers instead of through them.",
	"AVOID_FIRE":              "This monster will path around heat-related dangers instead of through them.",
	"AVOID_FALL":              "This monster will path around cliffs instead of off of them.",
	"PRIORITIZE_TARGETS":      "This monster will prioritize targets depending on their danger levels",
	"NOT_HALLU":               "Monsters that will NOT appear when player's producing hallucinations",
	"CANPLAY":                 "This creature can be played with if it's a pet.",
	"PET_MOUNTABLE":           "Creature can be ridden or attached to an harness.",
	"PET_HARNESSABLE":         "This monster can be harnessed when tamed.",
	"DOGFOOD":                 "Becomes friendly / tamed with dog food.",
	"MILKABLE":                "Produces milk when milked.",
	"SHEARABLE":               "This monster can be sheared for wool.",
	"NO_BREED":                "Creature doesn't reproduce even though it has reproduction data - useful when using copy-from to make child versions of adult creatures",
	"NO_FUNG_DMG":             "This monster can't be damaged by fungal spores and can't be fungalized either.",
	"PET_WONT_FOLLOW":         "This monster won't follow the player automatically when tamed.",
	"DRIPS_NAPALM":            "Occasionally drips napalm on move.",
	"DRIPS_GASOLINE":          "Occasionally drips gasoline on move.",
	"ELECTRIC_FIELD":          "This monster is surrounded by an electrical field that ignites flammable liquids near it",
	"LOUDMOVES":               "Makes move noises as if ~2 sizes louder even if flying.",
	"CAN_OPEN_DOORS":          "Can open doors on its path.",
	"STUN_IMMUNE":             "- This monster is immune to stun.",
	"DROPS_AMMO":              "This monster drops ammo. Should not be set for monsters that use pseudo ammo.",
	"INSECTICIDEPROOF":        "It's immune to insecticide even though it's made of bug flesh (\"iflesh\").",
	"RANGED_ATTACKER":         "Monster has any sort of ranged attack.",
	"CAMOUFLAGE":              "Stays invisible up to (current Perception + base Perception if the character has the Spotting proficiency) tiles away even in broad daylight. Monsters see it from the lower of vision_day and vision_night ranges.",
	"WATER_CAMOUFLAGE":        "If in water stays invisible up to (current Perception + base Perception if the character has the Spotting proficiency) tiles away even in broad daylight. Monsters see it from the lower of vision_day and vision_night ranges. Can also make it harder to see in deep water or across Z-levels if it is underwater and the viewer is not.",
	"MAX":                     "Sets the length of the flags - obviously must be LAST",
	"BIRDFOOD":                "Becomes friendly / tamed with bird food.",
	"CATFOOD":                 "Becomes friendly / tamed with cat food.",
	"CATTLEFODDER":            "Becomes friendly / tamed with cattle fodder.",
	"CURRENT":                 "this water is flowing.",
	"NOT_HALLUCINATION":       "This monster does not appear while the player is hallucinating.",
	"PET_HARNESSABLECreature": "can be attached to an harness.",
	"NULL":                    "Source use only."}
