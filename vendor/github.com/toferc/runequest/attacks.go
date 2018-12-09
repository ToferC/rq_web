package runequest

// Attack represents an offensive ability in Runequest
type Attack struct {
	Name             string
	Skill            *Skill
	Range            int
	StrikeRank       int
	BaseDamage       *DieCode
	AdditionalDamage []*DieCode
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
