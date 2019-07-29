package runequest

import (
	"fmt"
	"math"
)

// SkillCategory is a grouping of skills
type SkillCategory struct {
	Name  string
	Value int
}

func (sc SkillCategory) String() string {

	var text string

	if sc.Value > -1 {
		text += fmt.Sprintf("%s (+%d%%)", sc.Name, sc.Value)
	} else {
		text += fmt.Sprintf("%s (%d%%)", sc.Name, sc.Value)
	}
	return text
}

// CategoryOrder sets the order of skills in Runequest
var CategoryOrder = []string{"Agility", "Communication", "Knowledge", "Magic", "Manipulation", "Perception", "Stealth", "Melee",
	"Ranged", "Shield"}

type statMods struct {
	statistic string
	values    map[int]int
}

// DetermineSkillCategoryValues figures out base values for skill categories based on stats
func (c *Character) DetermineSkillCategoryValues() {

	for _, sc := range CategoryOrder {

		c.SkillCategories[sc] = &SkillCategory{}

		c.SkillCategories[sc].Name = sc
		c.SkillCategories[sc].Value = 0
	}

	for k, sc := range RQSkillCategories {
		// For each category

		for _, sm := range sc {
			// For each modifier in a category

			// Identify the stat
			stat := c.Statistics[sm.statistic]
			stat.UpdateStatistic()

			// Match to SkillCategory
			s := c.SkillCategories[k]

			// Map against specific values
			switch {
			case stat.Total == 0:
				s.Value = 0
			case stat.Total <= 4:
				s.Value += sm.values[4]
			case stat.Total <= 8:
				s.Value += sm.values[8]
			case stat.Total <= 12:
				s.Value += sm.values[12]
			case stat.Total <= 16:
				s.Value += sm.values[16]
			case stat.Total <= 20:
				s.Value += sm.values[20]
			case stat.Total > 20:
				if sm.values[20] > 0 {
					f := float64(stat.Total) - 20.0
					s.Value += sm.values[20] + int(math.Ceil(f/4))*5
				} else {
					f := float64(stat.Total) - 20.0
					s.Value += sm.values[20] - int(math.Ceil(f/4))*5
				}
			}
		}
	}

	for _, skill := range c.Skills {
		sc := c.SkillCategories[skill.Category]

		skill.CategoryValue = sc.Value
		skill.UpdateSkill()
	}
}

var minorPositive = map[int]int{
	4:  -5,
	8:  0,
	12: 0,
	16: 0,
	20: 5,
}

var majorPositive = map[int]int{
	4:  -10,
	8:  -5,
	12: 0,
	16: 5,
	20: 10,
}

var minorNegative = map[int]int{
	4:  5,
	8:  0,
	12: 0,
	16: 0,
	20: -5,
}

var majorNegative = map[int]int{
	4:  10,
	8:  5,
	12: 0,
	16: -5,
	20: -10,
}

// Common SkillCategory for combat skills & manipulation
var manipulationMods = []statMods{
	statMods{
		statistic: "STR",
		values:    minorPositive,
	},
	statMods{
		statistic: "DEX",
		values:    majorPositive,
	},
	statMods{
		statistic: "INT",
		values:    majorPositive,
	},
	statMods{
		statistic: "POW",
		values:    minorPositive,
	},
}

// RQSkillCategories is a map of skill categories
var RQSkillCategories = map[string][]statMods{
	// Agility
	"Agility": []statMods{
		statMods{
			statistic: "STR",
			values:    minorPositive,
		},
		statMods{
			statistic: "SIZ",
			values:    minorNegative,
		},
		statMods{
			statistic: "DEX",
			values:    majorPositive,
		},
		statMods{
			statistic: "POW",
			values:    minorPositive,
		},
	},
	// Communication
	"Communication": []statMods{
		statMods{
			statistic: "INT",
			values:    minorPositive,
		},
		statMods{
			statistic: "POW",
			values:    minorPositive,
		},
		statMods{
			statistic: "CHA",
			values:    majorPositive,
		},
	},
	// Knowledge
	"Knowledge": []statMods{
		statMods{
			statistic: "INT",
			values:    majorPositive,
		},
		statMods{
			statistic: "POW",
			values:    minorPositive,
		},
	},
	// Magic
	"Magic": []statMods{
		statMods{
			statistic: "POW",
			values:    majorPositive,
		},
		statMods{
			statistic: "CHA",
			values:    minorPositive,
		},
	},
	// Manipulation
	"Manipulation": manipulationMods,

	// Weapons
	"Melee": manipulationMods,

	"Ranged": manipulationMods,

	"Shield": manipulationMods,

	// Perception
	"Perception": []statMods{
		statMods{
			statistic: "INT",
			values:    majorPositive,
		},
		statMods{
			statistic: "POW",
			values:    minorPositive,
		},
	},
	// Stealth
	"Stealth": []statMods{
		statMods{
			statistic: "SIZ",
			values:    majorNegative,
		},
		statMods{
			statistic: "DEX",
			values:    majorPositive,
		},
		statMods{
			statistic: "INT",
			values:    majorPositive,
		},
		statMods{
			statistic: "POW",
			values:    minorNegative,
		},
	},
}
