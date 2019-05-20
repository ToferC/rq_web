package runequest

// NewCharacter generates a random starting character in Runequest
func NewCharacter(name string) *Character {
	c := Character{
		Name:           name,
		Role:           "Player Character",
		Statistics:     RuneQuestStats,
		Abilities:      Abilities,
		PowerRunes:     PowerRunes,
		ElementalRunes: ElementalRunes,
		ConditionRunes: ConditionRunes,
		RuneSpells:     map[string]*Spell{},
		SpiritMagic:    map[string]*Spell{},
		Homeland:       &Homeland{},
		Occupation:     &Occupation{},
		Cult:           &Cult{},
		HitLocations:   Locations,
		HitLocationMap: HitLocationMap,
		CreationSteps:  CreationStatus,

		MeleeAttacks:  map[string]*Attack{},
		RangedAttacks: map[string]*Attack{},

		SkillCategories: map[string]*SkillCategory{},
	}

	// Skills is a map of regular skills in Runequest
	c.Skills = map[string]*Skill{}

	//c.Skills["Dodge"].Base = c.Statistics["DEX"].Value * 2
	//c.Skills["Jump"].Base = c.Statistics["DEX"].Value * 3

	return &c
}
