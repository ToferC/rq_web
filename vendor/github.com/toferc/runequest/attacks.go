package runequest

import "fmt"

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
