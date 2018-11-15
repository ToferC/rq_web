package runequest

import (
	"fmt"
	"sort"
)

// Character represents a generic RPG character
type Character struct {
	Name             string
	Setting          string
	Description      string
	Race             *Race
	Homeland         *Homeland
	Occupation       *Occupation
	Cult             *Cult
	Clan             string
	Tribe            string
	Abilities        map[string]*Ability
	ElementalRunes   map[string]*Ability
	PowerRunes       map[string]*Ability
	Statistics       map[string]*Statistic
	Attributes       map[string]*Attribute
	DerivedMap       []string
	Skills           map[string]*Skill
	SkillMap         []string
	SkillCategories  map[string]*SkillCategory
	Advantages       map[string]*Advantage
	AdvantageMap     []string
	RuneSpells       map[string]*Spell
	SpiritMagic      map[string]*Spell
	Powers           map[string]*Power
	HitLocations     map[string]*HitLocation
	HitLocationMap   []string
	Equipment        []string
	Lunars           int
	Ransom           int
	StandardofLiving string
	PointCost        int
	InPlay           bool
	Updates          []*Update
}

// Update tracks live changes to Character
type Update struct {
	Date  string
	Event string
	Value int
}

// UpdateCharacter updates stats, runes and skills based on them
func (c *Character) UpdateCharacter() {
	c.TotalStatistics()
	c.DetermineSkillCategoryValues()

	// This can be optimized
	c.SetAttributes()
	c.UpdateAttributes()
}

func (c Character) String() string {
	text := c.Name
	text += fmt.Sprintf("\nHomeland: %s", c.Homeland.Name)
	text += fmt.Sprintf("\nOccupation: %s", c.Occupation.Name)
	text += fmt.Sprintf("\nCult: %s", c.Cult.Name)

	text += "\n\nStats:\n"
	for _, stat := range StatMap {
		text += fmt.Sprintf("%s\n", c.Statistics[stat])
	}

	text += "\nDerived Stats:\n"
	for _, ds := range c.Attributes {
		text += fmt.Sprintf("%s\n", ds)
	}

	text += "\nAbilities:"

	for _, at := range AbilityTypes {

		text += fmt.Sprintf("\n\n**%s**", at)

		for _, ability := range c.Abilities {

			if ability.Type == at {
				text += fmt.Sprintf("\n%s", ability)
			}
		}
	}

	text += "\nElemental Runes:"

	for _, ability := range c.ElementalRunes {

		text += fmt.Sprintf("\n%s", ability)
	}

	text += "\nPower Runes:"

	for _, ability := range c.ElementalRunes {

		text += fmt.Sprintf("\n%s", ability)
	}

	text += "\n\nSkills:"

	keys := make([]string, 0, len(c.Skills))
	for k := range c.Skills {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, co := range CategoryOrder {

		sc := c.SkillCategories[co]

		text += fmt.Sprintf("%s", sc)
		for _, skill := range keys {

			if c.Skills[skill].Category == sc.Name {
				text += fmt.Sprintf("\n%s", c.Skills[skill])
			}
		}

	}
	return text
}
