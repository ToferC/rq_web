package runequest

// Occupation represents a profession in Runequest
type Occupation struct {
	Name             string
	Description      string
	Notes            string
	Skills           []*Skill
	Weapons          []WeaponSelection
	SkillChoices     []SkillChoice
	StandardOfLiving string
	Income           int
	CultChoices      []Cult
	PassionList      []Ability
	Passions         []Ability
	Abilities        []Ability
	AbilityChoices   []AbilityChoice
	Ransom           int
	Equipment        []string
}

// WeaponSelection represents a Weapon Choice under Occupationg
type WeaponSelection struct {
	Description string
	Value       int
}

// WeaponCategories is an array of weapon choices for occupations
var WeaponCategories = []string{"Any", "Melee", "Ranged", "Shield", "Cultural"}

// Standards is an array with options for Standards of Living
var Standards = []string{"Destitute", "Poor", "Free", "Noble"}
