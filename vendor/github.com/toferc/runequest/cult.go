package runequest

// Cult represents a Religion in Runequest
type Cult struct {
	Name            string
	Description     string
	SubCult         bool
	Runes           []string
	Rank            string
	Skills          []*Skill
	SkillChoices    []SkillChoice
	Abilities       []Ability
	AbilityChoices  []AbilityChoice
	Weapons         []WeaponSelection
	PassionList     []Ability
	Passions        []Ability
	RuneSpells      []*Spell
	NumRunePoints   int
	NumSpiritMagic  int
	SpiritMagic     []*Spell
	ParentCult      *Cult
	AssociatedCults []Cult
}
