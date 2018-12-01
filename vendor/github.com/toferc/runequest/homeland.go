package runequest

// Homeland represents a homeland and cultural learnings
type Homeland struct {
	Name            string
	Description     string
	StatisticFrames map[string]*StatisticFrame
	Skills          []*Skill
	// Base Skill List
	SkillChoices []SkillChoice
	// Options for skills
	RuneBonus      string
	Abilities      []Ability
	AbilityChoices []AbilityChoice
	AbilityList    []Ability
	PassionList    []Ability
	Passions       []Ability
	Advantages     []Advantage
}

// StatisticFrame represents stat modifiers to a character
type StatisticFrame struct {
	Name     string
	Dice     int
	Modifier int
}
