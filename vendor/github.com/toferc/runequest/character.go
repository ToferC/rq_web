package runequest

import (
	"fmt"
	"sort"
)

// Character represents a generic RPG character
type Character struct {
	Name string
	Role string
	// Type of Character - Player, NPC, Creature, etc.
	Description string
	Race        *Race
	Homeland    *Homeland
	Occupation  *Occupation
	Cult        *Cult
	ExtraCults  []*ExtraCult
	Age         int
	Clan        string
	Tribe       string
	Abilities   map[string]*Ability
	// Passions and Reputation
	ElementalRunes map[string]*Ability
	// Elemental Runes
	PowerRunes       map[string]*Ability
	ConditionRunes   map[string]*Ability
	CoreRunes        []*Ability
	Statistics       map[string]*Statistic
	Attributes       map[string]*Attribute
	CurrentHP        int
	CurrentMP        int
	CurrentRP        int
	Movement         []*Movement
	DerivedMap       []string
	Skills           map[string]*Skill
	SkillMap         []string
	SkillCategories  map[string]*SkillCategory
	Advantages       map[string]*Advantage
	AdvantageMap     []string
	RuneSpells       map[string]*Spell
	SpiritMagic      map[string]*Spell
	Powers           map[string]*Power
	LocationForm     string
	HitLocations     map[string]*HitLocation
	HitLocationMap   []string
	MeleeAttacks     map[string]*Attack
	RangedAttacks    map[string]*Attack
	Equipment        []string
	Lunars           int
	Ransom           int
	StandardofLiving string
	InPlay           bool
	Updates          []*Update
	CreationSteps    map[string]bool

	Tags  []string
	Notes string
}

// CharacterRoles is an array of options for Character.Role
var CharacterRoles = []string{
	"Player Character",
	"Non-Player Character",
	"Animal",
	"Chaos",
	"Creature",
	"Demon",
	"Elemental",
	"Spirit",
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
	c.IDCoreRunes()

	if len(c.HitLocations) == 0 {
		hlForms := LocationForms[c.LocationForm]
		c.HitLocationMap = SortLocations(hlForms)
		c.HitLocations = hlForms
	}
}

// TotalStatistics updates values for stats after being modified
func (c *Character) TotalStatistics() {

	for _, s := range c.Statistics {

		s.UpdateStatistic()
	}
}

// CreationStatus tracks the completion of character creation
var CreationStatus = map[string]bool{
	"Base Choices":      false,
	"Personal History":  false,
	"Rune Affinities":   false,
	"Roll Stats":        false,
	"Apply Homeland":    false,
	"Apply Occupation":  false,
	"Apply Cult":        false,
	"Personal Skills":   false,
	"Finishing Touches": false,
	"Complete":          false,
}

func (c Character) String() string {
	text := c.Name
	text += fmt.Sprintf("\nHomeland: %s", c.Homeland.Name)
	text += fmt.Sprintf("\nOccupation: %s", c.Occupation.Name)
	text += fmt.Sprintf("\n%s of Cult: %s", c.Cult.Rank, c.Cult.Name)

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
