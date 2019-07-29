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
		LocationForm:   "Humanoids",
		HitLocations:   LocationForms["Humanoids"],
		HitLocationMap: SortLocations(LocationForms["Humanoid"]),
		CreationSteps:  CreationStatus,

		MeleeAttacks:  map[string]*Attack{},
		RangedAttacks: map[string]*Attack{},

		SkillCategories: map[string]*SkillCategory{},

		Movement: []*Movement{
			&Movement{
				Name:  "Ground",
				Value: 8,
			},
		},
	}

	// Skills is a map of regular skills in Runequest
	c.Skills = Skills

	return &c
}
