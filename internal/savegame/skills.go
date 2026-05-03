package savegame

// Skill and attribute identifiers as they appear in the save file.
// Order matches the original Delphi UI grouping so the Go UI can be driven
// from these tables.

var WeaponSkills = []string{
	"bluntWeapons",
	"smg",
	"brawling",
	"sniperRifle",
	"atWeapons",
	"bladedWeapons",
	"rifle",
	"energyWeapons",
	"shotgun",
	"handgun",
}

var GeneralSkills = []string{
	"calvinBackerSkill",
	"combatShooting",
	"outdoorsman",
	"bruteForce",
	"animalWhisperer",
	"spotLie",
	"intimidate",
	"perception",
	"leadership",
	"barter",
	"weaponSmith",
	"manipulate",
}

var TechnicalSkills = []string{
	"demolitions",
	"computerTech",
	"mechanicalRepair",
	"fieldMedic",
	"toasterRepair",
	"alarmDisarm",
	"doctor",
	"safecrack",
	"pickLock",
}

var Attributes = []string{
	"coordination",
	"luck",
	"awareness",
	"strength",
	"speed",
	"intelligence",
	"charisma",
}

// Human-readable labels keyed by the raw name used in the save file.
var SkillLabels = map[string]string{
	"bluntWeapons":      "Blunt Weapons",
	"smg":               "SMG",
	"brawling":          "Brawling",
	"sniperRifle":       "Sniper Rifle",
	"atWeapons":         "Big Guns",
	"bladedWeapons":     "Bladed Weapons",
	"rifle":             "Assault Rifle",
	"energyWeapons":     "Energy Weapons",
	"shotgun":           "Shotgun",
	"handgun":           "Handgun",
	"calvinBackerSkill": "Hard Ass",
	"combatShooting":    "Combat Shooting",
	"outdoorsman":       "Outdoorsman",
	"bruteForce":        "Brute Force",
	"animalWhisperer":   "Animal Whisperer",
	"spotLie":           "Kiss Ass",
	"intimidate":        "Intimidate",
	"perception":        "Perception",
	"leadership":        "Leadership",
	"barter":            "Barter",
	"weaponSmith":       "Weaponsmithing",
	"manipulate":        "Smart Ass",
	"demolitions":       "Demolitions",
	"computerTech":      "Computer Science",
	"mechanicalRepair":  "Mechanical Repair",
	"fieldMedic":        "Field Medic",
	"toasterRepair":     "Toaster Repair",
	"alarmDisarm":       "Alarm Disarm",
	"doctor":            "Surgeon",
	"safecrack":         "Safecracking",
	"pickLock":          "Lockpicking",
	"coordination":      "Coordination",
	"luck":              "Luck",
	"awareness":         "Awareness",
	"strength":          "Strength",
	"speed":             "Speed",
	"intelligence":      "Intelligence",
	"charisma":          "Charisma",
}

// Cumulative skill XP cost per skill level (index = level).
// e.g. level 5 costs 14 XP. Reaching level 10 costs 44.
var skillXpByLevel = [...]int{0, 2, 4, 6, 10, 14, 18, 24, 30, 36, 44}

// SkillLevel returns the displayed 0..10 level for a raw skillXp value.
// Values that don't match a level boundary round down.
func SkillLevel(rawXp int) int {
	level := 0
	for i, threshold := range skillXpByLevel {
		if rawXp >= threshold {
			level = i
		} else {
			break
		}
	}
	return level
}

// SkillXP returns the raw XP value to assign to reach a given 0..10 level.
// Out-of-range levels are clamped.
func SkillXP(level int) int {
	if level < 0 {
		level = 0
	}
	if level > 10 {
		level = 10
	}
	return skillXpByLevel[level]
}
