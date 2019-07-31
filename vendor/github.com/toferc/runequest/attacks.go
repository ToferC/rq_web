package runequest

import (
	"fmt"
	"strings"
)

// Attack represents an offensive ability in Runequest
type Attack struct {
	Name             string
	Skill            *Skill
	Weapon           *Weapon
	Range            int
	StrikeRank       int
	BaseDamage       *DieCode
	AdditionalDamage []*DieCode
	DamageString     string
	StrengthDamage   bool
	Special          string
}

// DieCode represents a single set of dice to roll 1d6+2 for example
type DieCode struct {
	Name     string
	NumDice  int
	DiceMax  int
	Modifier int
}

func (a *Attack) String() string {
	text := fmt.Sprintf("%s %d%% %s SR %d %d/%d HP %s",
		a.Weapon.Name,
		a.Skill.Total,
		a.DamageString,
		a.StrikeRank,
		a.Weapon.CurrentHP,
		a.Weapon.HP,
		a.Special)

	if a.Weapon.Range > 0 || a.Range > 0 {
		text += fmt.Sprintf("Rng %d", a.Weapon.Range)
	}
	return text
}

// UpdateAttacks reviews and updates weapon data on save
func (c *Character) UpdateAttacks() {

	db := c.Attributes["DB"]
	dbString := ""
	throwDB := ""

	if db.Text != "-" {
		dbString = db.Text

		if db.Base > 0 {
			throwDB = fmt.Sprintf("+%dD%d", db.Dice, db.Base/2)
		} else {
			throwDB = fmt.Sprintf("-%dD%d", db.Dice, db.Base/2)
		}
	}

	for _, m := range c.MeleeAttacks {
		m.Skill = c.Skills[m.Skill.Name]
		m.DamageString = m.Weapon.Damage + dbString
		m.StrikeRank = c.Attributes["DEXSR"].Base + c.Attributes["SIZSR"].Base + m.Weapon.SR
	}

	for _, r := range c.RangedAttacks {

		throw := false

		if strings.Contains(r.Weapon.Name, "Thrown") {
			throw = true
		}
		r.Weapon.Thrown = throw

		damage := ""

		if r.Weapon.Thrown {
			damage = r.Weapon.Damage + throwDB
		} else {
			damage = r.Weapon.Damage
		}

		r.Skill = c.Skills[r.Skill.Name]
		r.DamageString = damage
		r.StrikeRank = c.Attributes["DEXSR"].Base
	}
}
