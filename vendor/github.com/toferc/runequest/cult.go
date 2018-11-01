package runequest

// Cult represents a Religion in Runequest
type Cult struct {
	Name            string
	Description     string
	SubCult         bool
	Runes           []string
	Rank            string
	Skills          []Skill
	SkillChoices    []SkillChoice
	Abilities       []Ability
	AbilityChoices  []AbilityChoice
	Weapons         []WeaponSelection
	PassionList     []Ability
	Passions        []Ability
	RuneSpells      []Spell
	NumRunePoints   int
	NumSpiritMagic  int
	SpiritMagic     []Spell
	SubCults        []Cult
	AssociatedCults []Cult
}

// ChooseCult modifies a character's skills by Cult
func (c *Character) ChooseCult(hl *Cult) {

	if c.Cult == nil {
		// First Cult so apply all modifiers
		c.Cult = hl
		c.ApplyCult()
	} else {
		c.RemoveCult()
		c.Cult = hl
		c.ApplyCult()
		// Already has a Cult & need to remove previous Cult skills
	}
}

// ApplyCult applies a Cult Template to a character
func (c *Character) ApplyCult() {

	for _, s := range c.Cult.Skills {
		c.ModifySkill(s)
	}

	for _, choice := range c.Cult.SkillChoices {
		// Find number of skills
		l := len(choice.Skills)

		// Choose random index
		r := ChooseRandom(l)

		// Select index from choice.Skills
		selected := choice.Skills[r]
		c.Cult.Skills = append(c.Cult.Skills, selected)

		// Modify or add skill
		c.ModifySkill(selected)
	}

	passions := c.Cult.PassionList
	// Find number of abilities
	l := len(passions)

	// Choose random index
	r := ChooseRandom(l)

	// Select index from Passions
	selected := passions[r]
	c.Cult.Passions = append(c.Cult.Passions, selected)

	// Modify or add ability
	c.ModifyAbility(selected)

	// Same for abilities

	for _, choice := range c.Cult.AbilityChoices {
		// Find number of skills
		l = len(choice.Abilities)

		// Choose random index
		r := ChooseRandom(l)

		// Select index from choice.Abilities
		selected := choice.Abilities[r]
		c.Cult.Abilities = append(c.Cult.Abilities, selected)

		// Modify or add skill
		c.ModifyAbility(selected)
	}
}

// RemoveCult removes all Cult Modifers from a character
func (c *Character) RemoveCult() {

	for _, s := range c.Cult.Skills {
		s.HomelandValue = 0

		if s.Base > 0 {
			s.Base += s.Base * -1
		}

		c.ModifySkill(s)
	}

	for _, p := range c.Cult.Passions {
		p.HomelandValue = 0

		if p.Base > 0 {
			p.Base += p.Base * -1
		}

		c.ModifyAbility(p)
	}

	for _, a := range c.Cult.Abilities {
		a.HomelandValue = 0

		if a.Base > 0 {
			a.Base += a.Base * -1
		}

		c.ModifyAbility(a)
	}
}
