package runequest

// PersonalHistoryEvents is a map of events for character history
var PersonalHistoryEvents = map[string]Event{
	"1582_base": Event{
		Name:        "Base",
		Year:        1583,
		Start:       true,
		Description: "Lots of text",
		Participant: "Grandparent",
		HomelandModifiers: map[string]int{
			"Sartar":  -5,
			"Esrolia": 5,
		},
		Slug:           "1582_base",
		FollowingEvent: "1583_base",
		End:            false,
		Results: []EventResult{
			EventResult{
				Range:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				Description: "Go to war!",
				Skills: []Skill{
					Skill{
						Name: "Dodge",
						Updates: []*Update{
							&Update{
								Date:  "1583",
								Value: 10,
								Event: "1583 - go to war",
							},
						},
					},
				},
				Passions: []Ability{
					Ability{
						Name:       "Love (family)",
						CoreString: "Love",
						UserString: "family",
						Updates: []*Update{
							&Update{
								Date:  "1583",
								Value: 10,
								Event: "1583 - go to war",
							},
						},
					},
				},
				Lunars:               RollDice(10, 1, 0, 1) * 10,
				Reputation:           RollDice(4, 1, 0, 1),
				ImmediateFollowEvent: "1583_civil_war",
				ImmediateFollowMod:   0,
				NextFollowEvent:      "1583_civil_war",
				NextFollowMod:        0,
			},
			EventResult{
				Range:       []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
				Description: "Don't go to war",
				Skills: []Skill{
					Skill{
						Name: "Bargain",
						Updates: []*Update{
							&Update{
								Date:  "1624",
								Value: 10,
								Event: "1624 - don't go to war",
							},
						},
					},
				},
				Passions: []Ability{
					Ability{
						Name:       "Hate (Lunars)",
						CoreString: "Hate",
						UserString: "Lunars",
						Updates: []*Update{
							&Update{
								Date:  "1583",
								Value: 10,
								Event: "1583 - go to war",
							},
						},
					},
				},
				Lunars:          RollDice(10, 1, 0, 1) * 10,
				Reputation:      RollDice(4, 1, 0, 1),
				NextFollowEvent: "1583_civil_war",
				NextFollowMod:   0,
			},
		},
	},

	"1583_civil_war": Event{
		Name:        "Civil War in Esrolia",
		Year:        1583,
		Start:       false,
		Description: "Lots of text",
		Participant: "Grandparent",
		End:         true,
		HomelandModifiers: map[string]int{
			"Sartar":  0,
			"Esrolia": 0,
		},
		Slug: "1582_base",
		Results: []EventResult{
			EventResult{
				Range:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				Description: "Go to war!",
				Skills: []Skill{
					Skill{
						Name: "Dodge",
						Updates: []*Update{
							&Update{
								Date:  "1583",
								Value: 10,
								Event: "1583 - go to war",
							},
						},
					},
				},
				Passions: []Ability{
					Ability{
						Name:       "Love (family)",
						CoreString: "Love",
						UserString: "family",
						Updates: []*Update{
							&Update{
								Date:  "1583",
								Value: 10,
								Event: "1583 - go to war",
							},
						},
					},
				},
				Lunars:     RollDice(10, 1, 0, 1) * 10,
				Reputation: RollDice(4, 1, 0, 1),
			},
			EventResult{
				Range:       []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
				Description: "Don't go to war",
				Skills: []Skill{
					Skill{
						Name: "Bargain",
						Updates: []*Update{
							&Update{
								Date:  "1624",
								Value: 10,
								Event: "1624 - don't go to war",
							},
						},
					},
				},
				Passions: []Ability{
					Ability{
						Name:       "Hate (Lunars)",
						CoreString: "Hate",
						UserString: "Lunars",
						Updates: []*Update{
							&Update{
								Date:  "1583",
								Value: 10,
								Event: "1583 - go to war",
							},
						},
					},
				},
				Lunars:     RollDice(10, 1, 0, 1) * 10,
				Reputation: RollDice(4, 1, 0, 1),
			},
		},
	},
}
