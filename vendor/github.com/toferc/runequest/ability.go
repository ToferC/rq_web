package runequest

import (
	"fmt"
)

// Ability represents any non-Ability % ability in Runequest
type Ability struct {
	Name               string
	Type               string
	OpposedAbility     string
	UserChoice         bool
	CoreString         string
	UserString         string
	Base               int
	HomelandValue      int
	OccupationValue    int
	CultValue          int
	CreationBonusValue int
	Updates            []*Update
	Value              int
	InPlayXPValue      int
	Total              int
	Max                int
	Min                int
	ExperienceCheck    bool
}

// AbilityChoice is a choice between 2 or more skills
type AbilityChoice struct {
	Abilities []Ability
}

// UpdateAbility totals skill values based on input
func (a *Ability) UpdateAbility() {

	a.generateName()

	updates := 0

	for _, u := range a.Updates {
		updates += u.Value
	}

	a.Total = a.Base + a.HomelandValue + a.OccupationValue + a.CultValue + a.CreationBonusValue + a.InPlayXPValue + a.Value + updates
}

// generateName sets the ability map name
func (a *Ability) generateName() {

	var n string

	if a.UserString != "" {
		n = fmt.Sprintf("%s (%s)", a.CoreString, a.UserString)
	} else {
		n = a.CoreString
	}
	a.Name = n
}

// AbilityTypes is a slice of potential types for abilities
var AbilityTypes = []string{
	"Base",
	"Elemental Rune",
	"Power Rune",
	"Passion",
}

// ElementalRuneOrder sets the order for Elemental Runes
var ElementalRuneOrder = []string{
	"Air", "Darkness", "Earth", "Fire/Sky", "Moon", "Water",
}

// PowerRuneOrder sets the order for Power Runes
var PowerRuneOrder = []string{
	"Man", "Beast",
	"Fertility", "Death",
	"Harmony", "Disorder",
	"Truth", "Illusion",
	"Movement", "Stasis",
	"Chaos", "Dragonewt",
	"Plant", "Spirit",
	"Undead",
}

// ConditionRuneOrder sets the order for condition runes
var ConditionRuneOrder = []string{
	"Law", "Mastery", "Magic", "Infinity",
}

// PassionTypes represents different passions in Glorantha
var PassionTypes = []string{
	"Devotion", "Fear", "Hate", "Honor",
	"Loyalty", "Love",
}

func (a *Ability) String() string {

	a.UpdateAbility()

	var text = ""

	text += fmt.Sprintf("%s %d%%", a.Name, a.Total)

	return text
}

// ModifyElementalRune adds or modifies a Ability value
func (c *Character) ModifyElementalRune(a Ability) {

	a.generateName()

	// Modify existing Ability
	ab := c.ElementalRunes[a.Name]

	// Add or subtract s.Value from Ability
	ab.Value += a.Value
	ab.HomelandValue += a.HomelandValue
	ab.OccupationValue += a.OccupationValue
	ab.CultValue += a.CultValue
	ab.CreationBonusValue += a.CreationBonusValue
	ab.InPlayXPValue += a.InPlayXPValue

	a.UpdateAbility()
}

// ModifyPowerRune adds or modifies a Ability value
func (c *Character) ModifyPowerRune(a Ability) {

	a.generateName()

	// Modify existing Ability
	ab := c.PowerRunes[a.Name]

	// Add or subtract s.Value from Ability
	ab.Value += a.Value
	ab.HomelandValue += a.HomelandValue
	ab.OccupationValue += a.OccupationValue
	ab.CultValue += a.CultValue
	ab.CreationBonusValue += a.CreationBonusValue
	ab.InPlayXPValue += a.InPlayXPValue

	c.UpdateOpposedRune(ab)
}

// UpdateOpposedRune determines the opposing rune value
func (c *Character) UpdateOpposedRune(ab *Ability) {
	// Modify Opposing Rune if required
	if ab.OpposedAbility != "" {
		opposed := c.PowerRunes[ab.OpposedAbility]

		ab.UpdateAbility()

		// Maximum of 99 in a Power Rune
		if ab.Total > 99 {
			ab.Total = 99
		}

		opposed.UpdateAbility()

		diff := ab.Total + opposed.Total

		if diff > 100 {
			opposed.Base -= diff - 100
		}
	}
}

// ModifyAbility adds or modifies a Ability value
func (c *Character) ModifyAbility(a Ability) {

	a.generateName()

	if c.Abilities[a.Name] == nil {
		// New Ability
		c.Abilities[a.Name] = &Ability{
			Name:       a.Name,
			CoreString: a.CoreString,
			UserString: a.UserString,
			Base:       a.Base,
			Type:       a.Type,
			Updates:    []*Update{},
		}
	} else {
		// Modify existing Ability
		ab := c.Abilities[a.Name]

		// Add or subtract s.Value from Ability
		if a.HomelandValue > 0 {
			ab.HomelandValue = a.HomelandValue
		}

		if a.OccupationValue > 0 {
			ab.OccupationValue = a.OccupationValue
		}

		if a.CultValue > 0 {
			ab.CultValue = a.CultValue
		}

		if a.CreationBonusValue > 0 {
			ab.CreationBonusValue = a.CreationBonusValue
		}

		if a.InPlayXPValue > 0 {
			ab.InPlayXPValue = a.InPlayXPValue
		}
	}
	a.UpdateAbility()
}

// Abilities is a map of the basic abilities in Runequest
var Abilities = map[string]*Ability{
	"Reputation": &Ability{
		CoreString: "Reputation",
		Type:       "Base",
		Value:      5,
	},
}

// ElementalRunes is a map of Rune abilities
var ElementalRunes = map[string]*Ability{
	// Elemental Runes
	"Darkness": &Ability{
		CoreString: "Darkness",
		Type:       "Elemental Rune",
		Value:      0,
	},
	"Air": &Ability{
		CoreString: "Air",
		Type:       "Elemental Rune",
		Value:      0,
	},
	"Water": &Ability{
		CoreString: "Water",
		Type:       "Elemental Rune",
		Value:      0,
	},
	"Earth": &Ability{
		CoreString: "Earth",
		Type:       "Elemental Rune",
		Value:      0,
	},
	"Fire/Sky": &Ability{
		CoreString: "Fire/Sky",
		Type:       "Elemental Rune",
		Value:      0,
	},
	"Moon": &Ability{
		CoreString: "Moon",
		Type:       "Elemental Rune",
		Value:      0,
	},
}

// PowerRunes is a map of Power Runes
var PowerRunes = map[string]*Ability{
	// Power Runes
	"Man": &Ability{
		CoreString:     "Man",
		Type:           "Form Rune",
		OpposedAbility: "Beast",
		Base:           50,
	},
	"Beast": &Ability{
		CoreString:     "Beast",
		Type:           "Form Rune",
		OpposedAbility: "Man",
		Base:           50,
	},
	"Fertility": &Ability{
		CoreString:     "Fertility",
		Type:           "Power Rune",
		OpposedAbility: "Death",
		Base:           50,
	},
	"Death": &Ability{
		CoreString:     "Death",
		Type:           "Power Rune",
		OpposedAbility: "Fertility",
		Base:           50,
	},
	"Harmony": &Ability{
		CoreString:     "Harmony",
		Type:           "Power Rune",
		OpposedAbility: "Disorder",
		Base:           50,
	},
	"Disorder": &Ability{
		CoreString:     "Disorder",
		Type:           "Power Rune",
		OpposedAbility: "Harmony",
		Base:           50,
	},
	"Truth": &Ability{
		CoreString:     "Truth",
		Type:           "Power Rune",
		OpposedAbility: "Illusion",
		Base:           50,
	},
	"Illusion": &Ability{
		CoreString:     "Illusion",
		Type:           "Power Rune",
		OpposedAbility: "Truth",
		Base:           50,
	},
	"Stasis": &Ability{
		CoreString:     "Stasis",
		Type:           "Power Rune",
		OpposedAbility: "Movement",
		Base:           50,
	},
	"Movement": &Ability{
		CoreString:     "Movement",
		Type:           "Power Rune",
		OpposedAbility: "Stasis",
		Base:           50,
	},
	"Chaos": &Ability{
		CoreString: "Chaos",
		Type:       "Form Rune",
		Base:       0,
	},
	"Dragonewt": &Ability{
		CoreString:     "Dragonewt",
		Type:           "Form Rune",
		OpposedAbility: "Man",
		Base:           0,
	},
	"Plant": &Ability{
		CoreString:     "Plant",
		Type:           "Form Rune",
		OpposedAbility: "Man",
		Base:           0,
	},
	"Spirit": &Ability{
		CoreString:     "Spirit",
		Type:           "Form Rune",
		OpposedAbility: "Man",
		Base:           0,
	},
	"Undead": &Ability{
		CoreString: "Undead",
		Type:       "Form Rune",
		Base:       0,
	},
}

// ConditionRunes is a map of Power Runes
var ConditionRunes = map[string]*Ability{
	// Condition Runes
	"Law": &Ability{
		CoreString: "Law",
		Type:       "Condition Rune",
		Base:       0,
	},
	"Mastery": &Ability{
		CoreString: "Mastery",
		Type:       "Condition Rune",
		Base:       0,
	},
	"Magic": &Ability{
		CoreString: "Magic",
		Type:       "Condition Rune",
		Base:       0,
	},
	"Infinity": &Ability{
		CoreString: "Infinity",
		Type:       "Condition Rune",
		Base:       0,
	},
}
