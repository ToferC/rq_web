package runequest

import (
	"fmt"
	"sort"
)

// Creature represents a generic RPG Creature
type Creature struct {
	Name string
	Role string
	// Type of Creature
	Description string
	
	Age         int
	Affiliations        string
	Abilities   map[string]*Ability
	// Passions and Reputation
	ElementalRunes map[string]*Ability
	// Elemental Runes
	PowerRunes       map[string]*Ability
	StatisticFrames  map[string]*StatisticFrame
	Statistics       map[string]*Statistic
	Attributes       map[string]*Attribute
	CurrentHP        int
	CurrentMP        int
	CurrentRP        int
	Cults map[string]int
	Move             int
	DerivedMap       []string
	Skills           map[string]*Skill
	SkillMap         []string
	SkillCategories  map[string]*SkillCategory
	RuneSpells       map[string]*Spell
	SpiritMagic      map[string]*Spell
	Powers           map[string]*Power
	HitLocations     map[string]*HitLocation
	HitLocationMap   []string
	MeleeAttacks     map[string]*Attack
	RangedAttacks    map[string]*Attack
	Equipment        []string
}

// CreatureRoles is an array of options for Creature.Role
var CreatureRoles = []string{
	"Non-Player Character",
	"Animal",
	"Chaos",
	"Creature",
	"Demon",
	"Elemental",
	"Spirit",
}

// TotalStatistics updates values for stats after being modified
func (c *Creature) TotalStatistics() {

	for _, s := range c.Statistics {

		s.UpdateStatistic()
	}
}

// UpdateCreature updates stats, runes and skills based on them
func (c *Creature) UpdateCreature() {
	c.TotalStatistics()
	c.DetermineSkillCategoryValues()

	// This can be optimized
	c.SetAttributes()
	c.UpdateAttributes()

	if len(c.HitLocations) == 0 {
		c.HitLocationMap = HitLocationMap
		c.HitLocations = Locations
	}
}

// DetermineSkillCategoryValues figures out base values for skill categories based on stats
func (c *Creature) DetermineSkillCategoryValues() {

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
					s.Value += sm.values[20] + ((stat.Total-20)/4)*5
				} else {
					s.Value += sm.values[20] - ((stat.Total-20)/4)*5
				}
			}
		}
	}

func (c Creature) String() string {
	
	text := c.Name
	text += fmt.Sprintf("\nAffiliations: %s", c.Affiliations)

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
