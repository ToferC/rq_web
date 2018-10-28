package runequest

// Homelands is a map of possible homelands in Runequest
var Homelands = map[string]*Homeland{
	// Sartar
	"Sartar": &Homeland{
		Name:      "Sartar",
		RuneBonus: "Air",
		Skills: []Skill{
			Skill{
				Name:          "Ride",
				CoreString:    "Ride",
				UserChoice:    true,
				UserString:    "Horse",
				Base:          5,
				HomelandValue: 5,
				Category:      "Agility",
			},
			Skill{
				CoreString:    "Dance",
				HomelandValue: 5,
			},
			Skill{
				CoreString:    "Sing",
				HomelandValue: 10,
			},
			Skill{
				Name:       "Speak",
				CoreString: "Speak",
				UserString: "Heortling",
				Base:       50,
				Category:   "Communication",
			},
			Skill{
				Name:          "Speak",
				UserString:    "Tradetalk",
				HomelandValue: 10,
				Category:      "Communication",
			},
			Skill{
				Name:       "Customs",
				CoreString: "Customs",
				UserString: "Heortling",
				Base:       25,
				Category:   "Knowledge",
			},
			Skill{
				CoreString:    "Farm",
				HomelandValue: 20,
			},
			Skill{
				CoreString:    "Herd",
				HomelandValue: 10,
			},
			Skill{
				CoreString:    "Spirit Combat",
				HomelandValue: 15,
			},
			Skill{
				CoreString:    "Dagger",
				HomelandValue: 10,
			},
			Skill{
				CoreString:    "Broadsword",
				HomelandValue: 15,
			},
			Skill{
				CoreString:    "Battle Axe",
				HomelandValue: 10,
			},
			Skill{
				CoreString:    "1H Spear",
				HomelandValue: 10,
			},
			Skill{
				CoreString:    "Javelin",
				HomelandValue: 10,
			},
			Skill{
				CoreString:    "Medium Shield",
				HomelandValue: 15,
			},
			Skill{
				CoreString:    "Large Shield",
				HomelandValue: 10,
			},
		},
		// Skill Choices
		SkillChoices: []SkillChoice{
			// Choice of 2 skills
			SkillChoice{
				Skills: []Skill{
					// Skill 1
					Skill{
						CoreString:    "Composite Bow",
						Base:          5,
						HomelandValue: 10,
					},
					// Skill 2
					Skill{
						CoreString:    "Sling",
						HomelandValue: 10,
					},
				},
			},
		},
		// Rune Affinities
		AbilityList: []Ability{
			// Ability 1
			Ability{
				Name:          "Air",
				CoreString:    "Air",
				Type:          "Elemental Rune",
				HomelandValue: 10,
			},
		},
		// Passions
		PassionList: []Ability{
			// Ability 1
			Ability{
				Name:          "Loyalty",
				CoreString:    "Loyalty",
				UserString:    "Sartar",
				UserChoice:    true,
				Type:          "Passion",
				Base:          60,
				HomelandValue: 10,
			},
			// Ability 2
			Ability{
				Name:          "Loyalty",
				CoreString:    "Loyalty",
				UserString:    "clan",
				Type:          "Passion",
				UserChoice:    true,
				Base:          60,
				HomelandValue: 10,
			},
			// Ability 3
			Ability{
				Name:          "Loyalty",
				CoreString:    "Loyalty",
				UserString:    "tribe",
				Type:          "Passion",
				UserChoice:    true,
				Base:          60,
				HomelandValue: 10,
			},
		},
	},
	// Esrolia
}

// Occupations is a map of possible Occupations in Runequest
var Occupations = map[string]*Occupation{
	"Farmer": &Occupation{
		Name: "Farmer",
		Skills: []Skill{
			Skill{
				Name:            "Occupation Lore",
				CoreString:      "Occupation Lore",
				UserString:      "Local",
				UserChoice:      true,
				OccupationValue: 15,
				Category:        "Knowledge",
			},
			Skill{
				CoreString:      "Farm",
				OccupationValue: 30,
			},
			Skill{
				Name:            "Craft",
				UserChoice:      true,
				CoreString:      "Craft",
				UserString:      "Arms",
				Base:            10,
				OccupationValue: 15,
				Category:        "Manipulation",
			},
			Skill{
				CoreString:      "First Aid",
				OccupationValue: 10,
			},
			Skill{
				CoreString:      "Scan",
				OccupationValue: 10,
			},
			Skill{
				CoreString:      "Herd",
				OccupationValue: 15,
			},
			Skill{
				CoreString:      "Manage Household",
				OccupationValue: 30,
			},
			Skill{
				CoreString:      "Medium Shield",
				OccupationValue: 15,
			},
			Skill{
				CoreString:      "Broadsword",
				OccupationValue: 15,
			},
		},
		// Skill Choices
		SkillChoices: []SkillChoice{
			// Choice of 2 skills
			SkillChoice{
				Skills: []Skill{
					// Skill 1
					Skill{
						Name:            "Jump",
						CoreString:      "Jump",
						OccupationValue: 10,
					},
					// Skill 2
					Skill{
						Name:            "Climb",
						CoreString:      "Climb",
						OccupationValue: 10,
					},
				},
			},
		},
		// Passions
		PassionList: []Ability{
			// Ability 1
			Ability{
				Name:            "Love",
				CoreString:      "Love",
				UserString:      "family",
				UserChoice:      true,
				Type:            "Passion",
				Base:            60,
				OccupationValue: 10,
			},
			// Ability 2
			Ability{
				Name:            "Loyalty",
				CoreString:      "Loyalty",
				UserString:      "clan",
				Type:            "Passion",
				UserChoice:      true,
				Base:            60,
				OccupationValue: 10,
			},
			// Ability 3
			Ability{
				Name:            "Loyalty",
				CoreString:      "Loyalty",
				UserString:      "tribe",
				Type:            "Passion",
				UserChoice:      true,
				Base:            60,
				OccupationValue: 10,
			},
		},
	},
	// Esrolia
}

// Cults is a map of possible Cults in Runequest
var Cults = map[string]*Cult{
	"Orlanth": &Cult{
		Name: "Orlanth",
		Skills: []Skill{
			Skill{
				CoreString: "Cult Lore",
				CultValue:  15,
				Category:   "Knowledge",
			},
			Skill{
				CoreString: "Worship",
				CultValue:  15,
				Category:   "Magic",
			},
			Skill{
				CoreString: "Meditate",
				CultValue:  25,
			},
			Skill{
				CoreString: "Orate",
				CultValue:  30,
			},
			Skill{
				CoreString: "Sing",
				CultValue:  10,
			},
			Skill{
				CoreString: "Speak",
				UserString: "Stormspeech",
				UserChoice: true,
				CultValue:  10,
				Category:   "Communication",
			},
		},
		// Passions
		PassionList: []Ability{
			// Ability 1
			Ability{
				Name:       "Devotion",
				CoreString: "Devotion",
				UserString: "Orlanth",
				UserChoice: true,
				Type:       "Passion",
				Base:       60,
				CultValue:  10,
			},
			// Ability 2
			Ability{
				Name:       "Hate",
				CoreString: "Hate",
				UserString: "Chaos",
				Type:       "Passion",
				UserChoice: true,
				Base:       60,
				CultValue:  10,
			},
			// Ability 3
			Ability{
				Name:       "Loyalty",
				CoreString: "Loyalty",
				UserString: "temple",
				Type:       "Passion",
				UserChoice: true,
				Base:       60,
				CultValue:  10,
			},
			// Ability 4
			Ability{
				Name:       "Honor",
				CoreString: "Honor",
				Type:       "Passion",
				Base:       60,
				CultValue:  10,
			},
		},
	},
	// Esrolia
}
