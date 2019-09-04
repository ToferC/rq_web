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
	BoundSpirits     []*BoundSpirit
	Income           int
	Lunars           int
	Ransom           int
	StandardofLiving string
	InPlay           bool
	Updates          []*Update
	CreationSteps    map[string]bool

	History []*Event

	Grandparent *FamilyMember
	Parent      *FamilyMember

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
	c.UpdateAttacks()

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
	text := fmt.Sprintf("\n--- %s ---\n", c.Name)

	if len(c.CoreRunes) > 0 {
		text += "Runes: "
		for _, cr := range c.CoreRunes {
			text += fmt.Sprintf("%s ", cr.Name)
		}
		text += "\n"
	}

	text += fmt.Sprintf("\nType: %s", c.Role)

	if c.Role == "Player Character" {
		text += fmt.Sprintf("\nHomeland: %s", c.Homeland.Name)
		text += fmt.Sprintf("\nOccupation: %s", c.Occupation.Name)
		text += fmt.Sprintf("\n%s of Cult: %s", c.Cult.Rank, c.Cult.Name)

		text += fmt.Sprintf("\nStandard of Living: %s", c.StandardofLiving)
		text += fmt.Sprintf("\nIncome: %d L", c.Income)
		text += fmt.Sprintf("\nRansom: %d L", c.Ransom)
	}

	text += fmt.Sprintf("\n\nDescription:\n%s", c.Description)

	if len(c.Statistics) > 0 {
		text += "\n\nStats:\n"
		for _, stat := range StatMap {
			if c.Statistics[stat].Total > 0 {
				text += fmt.Sprintf("%s (%d%%)\n", c.Statistics[stat],
					c.Statistics[stat].Total*5)
			}
		}
	}

	if len(c.Attributes) > 0 {
		text += "\nDerived Stats:\n"
		for _, ds := range c.Attributes {
			text += fmt.Sprintf("%s\n", ds)
		}
	}

	if len(c.Movement) > 0 {
		text += "\nMovement:\n"
		for _, m := range c.Movement {
			text += fmt.Sprintf("%s\n", m)
		}
	}

	if len(c.Abilities) > 0 {
		text += "\nPassions & Reputations:"

		for _, ability := range c.Abilities {
			text += fmt.Sprintf("\n%s", ability)
		}
	}

	if c.Cult.Name != "" {
		text += fmt.Sprintf("\n\n**Cults:\n%s - %s - Rune Points: %d", c.Cult.Name, c.Cult.Rank, c.Cult.NumRunePoints)
	}

	if len(c.ExtraCults) > 0 {
		for _, ec := range c.ExtraCults {
			text += fmt.Sprintf("\n%s - %s Rune Points: %d", ec.Name, ec.Rank, ec.RunePoints)
		}
	}

	if len(c.ElementalRunes) > 0 {
		text += "\n\nElemental Runes:"

		for _, ability := range c.ElementalRunes {
			if ability.Total != 0 {
				text += fmt.Sprintf("\n%s", ability)
			}
		}
	}

	if len(c.PowerRunes) > 0 {
		text += "\n\nPower Runes:"

		for _, ability := range c.PowerRunes {
			if ability.Total != 0 {
				text += fmt.Sprintf("\n%s", ability)
			}
		}
	}

	if len(ConditionRunes) > 0 {
		text += "\n\nCondition Runes:"

		for _, ability := range c.ConditionRunes {
			if ability.Total != 0 {
				text += fmt.Sprintf("\n%s", ability)
			}
		}
	}

	if len(c.Skills) > 0 {
		text += "\n\nSkills:"

		if len(c.Skills) > 20 {

			keys := make([]string, 0, len(c.Skills))
			for k := range c.Skills {
				keys = append(keys, k)
			}

			sort.Strings(keys)

			for _, co := range CategoryOrder {

				sc := c.SkillCategories[co]

				text += fmt.Sprintf("\n**%s**\n", sc)
				for _, skill := range keys {

					if c.Skills[skill].Category == sc.Name {
						text += fmt.Sprintf("%s\n", c.Skills[skill])
					}
				}
			}
		} else {
			text += "\n"
			for _, skill := range c.Skills {
				text += fmt.Sprintf("%s\n", skill)
			}
		}
	}

	if len(c.SpiritMagic) > 0 {
		text += "\nSpirit Magic:\n"
		for _, sm := range c.SpiritMagic {
			text += fmt.Sprintf("%s\n", sm)
		}
	}

	if len(c.RuneSpells) > 0 {
		text += "\nRune Spells:\n"
		for _, rs := range c.RuneSpells {
			text += fmt.Sprintf("%s\n", rs)
		}
	}

	if len(c.Powers) > 0 {
		text += "\nPowers:\n"
		for _, p := range c.Powers {
			text += fmt.Sprintf("%s: %s\n", p.Name, p.Description)
		}
	}

	if len(c.BoundSpirits) > 0 {
		text += "\nSpirits & Matrices:\n"
		for _, bs := range c.BoundSpirits {
			text += bs.String() + "\n"
		}
	}

	if len(c.MeleeAttacks) > 0 {
		text += "\nMelee Attacks:\n"
		for _, m := range c.MeleeAttacks {
			text += fmt.Sprintf("%s\n", m)
		}
	}

	if len(c.RangedAttacks) > 0 {
		text += "\nRanged Attacks:\n"
		for _, r := range c.RangedAttacks {
			text += fmt.Sprintf("%s\n", r)
		}
	}

	if len(c.HitLocations) > 0 {
		text += "\nHit Locations:\n"
		for _, hlm := range c.HitLocationMap {
			for k, v := range c.HitLocations {
				if k == hlm {
					text += fmt.Sprintf("%s", v)
				}
			}
		}
	}

	if len(c.Equipment) > 0 {
		text += "\nEquipment:\n"
		for _, e := range c.Equipment {
			text += fmt.Sprintf("%s\n", e)
		}
	}

	return text
}
