package oneroll

// NewSRCharacter generates an ORE WT character
func NewSRCharacter(name string) *Character {

	c := Character{
		Name: name,
	}

	c.Setting = "SR"

	c.Archetype = new(Archetype)

	// WTStats sets the order for Character.Statistics
	c.StatMap = []string{"Body", "Coordination", "Sense", "Mind", "Command", "Charm"}

	c.Statistics = map[string]*Statistic{
		"Body": &Statistic{
			Name: "Body",
			Dice: &DiePool{
				Normal:  2,
				Hard:    0,
				GoFirst: 0,
			},
		},
		"Coordination": &Statistic{
			Name: "Coordination",
			Dice: &DiePool{
				Normal: 2,
			},
		},
		"Sense": &Statistic{
			Name: "Sense",
			Dice: &DiePool{
				Normal: 2,
			},
		},
		"Mind": &Statistic{
			Name: "Mind",
			Dice: &DiePool{
				Normal: 2,
			},
		},
		"Command": &Statistic{
			Name: "Command",
			Dice: &DiePool{
				Normal: 2,
			},
			EffectsWill: true,
		},
		"Charm": &Statistic{
			Name: "Charm",
			Dice: &DiePool{
				Normal: 2,
			},
			EffectsWill: true,
		},
	}

	// Declare stat pointers

	body := c.Statistics["Body"]
	coordination := c.Statistics["Coordination"]
	sense := c.Statistics["Sense"]
	mind := c.Statistics["Mind"]
	command := c.Statistics["Command"]
	charm := c.Statistics["Charm"]

	c.LocationMap = []string{"Head", "Body", "Left Arm", "Right Arm",
		"Left Leg", "Right Leg"}

	c.HitLocations = map[string]*Location{
		"Head": &Location{
			Name:     "Head",
			HitLoc:   []int{10},
			Boxes:    4,
			Shock:    []bool{},
			Kill:     []bool{},
			LAR:      0,
			HAR:      0,
			Disabled: false,
		},
		"Body": &Location{
			Name:     "Body",
			HitLoc:   []int{7, 8, 9},
			Boxes:    10,
			Shock:    []bool{},
			Kill:     []bool{},
			LAR:      0,
			HAR:      0,
			Disabled: false,
		},
		"Left Arm": &Location{
			Name:     "Left Arm",
			HitLoc:   []int{5, 6},
			Boxes:    6,
			Shock:    []bool{},
			Kill:     []bool{},
			LAR:      0,
			HAR:      0,
			Disabled: false,
		},
		"Right Arm": &Location{
			Name:     "Right Arm",
			HitLoc:   []int{3, 4},
			Boxes:    6,
			Shock:    []bool{},
			Kill:     []bool{},
			LAR:      0,
			HAR:      0,
			Disabled: false,
		},
		"Left Leg": &Location{
			Name:     "Left Leg",
			HitLoc:   []int{2},
			Boxes:    6,
			Shock:    []bool{},
			Kill:     []bool{},
			LAR:      0,
			HAR:      0,
			Disabled: false,
		},
		"Right Leg": &Location{
			Name:     "Right Leg",
			HitLoc:   []int{1},
			Boxes:    6,
			Shock:    []bool{},
			Kill:     []bool{},
			LAR:      0,
			HAR:      0,
			Disabled: false,
		},
	}

	for _, v := range c.HitLocations {
		v.FillWounds()
	}

	// Shadowrun skills
	c.Skills = map[string]*Skill{
		// Body Skills
		"Athletics": &Skill{
			Name: "Athletics",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: body,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
				Wiggle: 0,
			},
		},
		"Block": &Skill{
			Name: "Block",
			Quality: &Quality{
				Type:  "Defend",
				Level: 0,
			},
			LinkStat: body,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
				Wiggle: 0,
			},
		},
		"Brawling": &Skill{
			Name: "Brawling",
			Quality: &Quality{
				Type:  "Attack",
				Level: 0,
			},
			LinkStat: body,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
				Wiggle: 0,
			},
		},
		"Endurance": &Skill{
			Name: "Endurance",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: body,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
				Wiggle: 0,
			},
		},
		// Coordination Skills
		"Acrobatics": &Skill{
			Name: "Acrobatics",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
				Wiggle: 0,
			},
		},
		"Close Combat": &Skill{
			Name: "Close Combat",
			Quality: &Quality{
				Type:  "Attack",
				Level: 0,
			},
			LinkStat: coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
			},
		},
		"Dodge": &Skill{
			Name: "Dodge",
			Quality: &Quality{
				Type:  "Defend",
				Level: 0,
			},
			LinkStat: coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
			},
		},
		"Escape Artist": &Skill{
			Name: "Escape Artist",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
			},
		},
		"Pilot": &Skill{
			Name: "Pilot",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
			},
			ReqSpec:        true,
			Specialization: "Ground",
		},
		"Small Arms": &Skill{
			Name: "Small Arms",
			Quality: &Quality{
				Type:  "Attack",
				Level: 0,
			},
			LinkStat: coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
			},
		},
		"Stealth": &Skill{
			Name: "Stealth",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
			},
		},
		"Heavy Weapons": &Skill{
			Name: "Heavy Weapons",
			Quality: &Quality{
				Type:  "Attack",
				Level: 0,
			},
			LinkStat: coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
			},
		},
		"Vehicle Weapons": &Skill{
			Name: "Vehicle Weapons",
			Quality: &Quality{
				Type:  "Attack",
				Level: 0,
			},
			LinkStat: coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
			},
		},
		// Sense Skills
		"Artisan": &Skill{
			Name: "Artisan",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: sense,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Empathy": &Skill{
			Name: "Empathy",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: sense,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Perception": &Skill{
			Name: "Perception",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: sense,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Scrutiny": &Skill{
			Name: "Scrutiny",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: sense,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Tracking": &Skill{
			Name: "Tracking",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: sense,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Disguise": &Skill{
			Name: "Disguise",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: sense,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		// Mind Skills
		"Computer": &Skill{
			Name: "Computer",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Armorer": &Skill{
			Name: "Armorer",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Cybertechnology": &Skill{
			Name: "Cybertechnology",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"First Aid": &Skill{
			Name: "First Aid",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Knowledge": &Skill{
			Name: "Knowledge",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: mind,
			Dice: &DiePool{
				Normal: 0,
			},
			ReqSpec:        true,
			Specialization: "Biology",
			Free:           true,
		},
		"Engineering": &Skill{
			Name: "Engineering",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Demolitions": &Skill{
			Name: "Demolitions",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Electronic Warfare": &Skill{
			Name: "Electronic Warfare",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: mind,
			Dice: &DiePool{
				Normal: 0,
			},
		}, "Hacking": &Skill{
			Name: "Hacking",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Languages": &Skill{
			Name: "Languages",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: mind,
			Dice: &DiePool{
				Normal: 0,
			},
			ReqSpec:        true,
			Free:           true,
			Specialization: "Chinese",
		},
		"Medicine": &Skill{
			Name: "Medicine",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Navigation": &Skill{
			Name: "Navigation",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Data Search": &Skill{
			Name: "Data Search",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Security Systems": &Skill{
			Name: "Security Systems",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Software": &Skill{
			Name: "Software",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Streetwise": &Skill{
			Name: "Streetwise",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Tactics": &Skill{
			Name: "Tactics",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: mind,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		// Charm Skills
		"Con": &Skill{
			Name: "Con",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: charm,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Ettiquette": &Skill{
			Name: "Ettiquette",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: charm,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Lie": &Skill{
			Name: "Lie",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: charm,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Performance": &Skill{
			Name: "Performance",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: charm,
			Dice: &DiePool{
				Normal: 0,
			},
			ReqSpec:        true,
			Specialization: "Standup",
		},
		"Persuasion": &Skill{
			Name: "Persuasion",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: charm,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Wealth": &Skill{
			Name: "Wealth",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: charm,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		// Command Skills
		"Arcane": &Skill{
			Name: "Arcane",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: command,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Interrogation": &Skill{
			Name: "Interrogation",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: command,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Intimidation": &Skill{
			Name: "Intimidation",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: command,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Leadership": &Skill{
			Name: "Leadership",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: command,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Stability": &Skill{
			Name: "Stability",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: command,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Survival": &Skill{
			Name: "Survival",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: command,
			Dice: &DiePool{
				Normal: 0,
			},
		},
	}

	c.Powers = map[string]*Power{}

	return &c
}
