package oneroll

// NewReignCharacter generates an ORE WT character
func NewReignCharacter(name string) *Character {

	c := Character{
		Name: name,
	}

	c.Setting = "RE"

	// WTStats sets the order for Character.Statistics
	c.StatMap = []string{"Body", "Coordination", "Sense", "Knowledge", "Command", "Charm"}

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
		"Knowledge": &Statistic{
			Name: "Knowledge",
			Dice: &DiePool{
				Normal: 2,
			},
		},
		"Command": &Statistic{
			Name: "Command",
			Dice: &DiePool{
				Normal: 2,
			},
		},
		"Charm": &Statistic{
			Name: "Charm",
			Dice: &DiePool{
				Normal: 2,
			},
		},
	}

	// Declare stat pointers

	body := c.Statistics["Body"]
	coordination := c.Statistics["Coordination"]
	sense := c.Statistics["Sense"]
	knowledge := c.Statistics["Knowledge"]
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
			},
		},
		"Parry": &Skill{
			Name: "Parry",
			Quality: &Quality{
				Type:  "Defend",
				Level: 0,
			},
			LinkStat: body,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Fight": &Skill{
			Name: "Fight",
			Quality: &Quality{
				Type:  "Attack",
				Level: 0,
			},
			LinkStat: body,
			Dice: &DiePool{
				Normal: 0,
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
			},
		},
		"Vigor": &Skill{
			Name: "Vigor",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: body,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Run": &Skill{
			Name: "Run",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: body,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		// Coordination Skills
		"Climb": &Skill{
			Name: "Climb",
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
		"Perform": &Skill{
			Name: "Perform",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: coordination,
			Dice: &DiePool{
				Normal: 0,
			},
			ReqSpec:        true,
			Specialization: "Juggler",
		},
		"Ride": &Skill{
			Name: "Ride",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: coordination,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Weapon": &Skill{
			Name: "Weapon",
			Quality: &Quality{
				Type:  "Attack",
				Level: 0,
			},
			LinkStat: coordination,
			Dice: &DiePool{
				Normal: 0,
				Hard:   0,
			},
			ReqSpec:        true,
			Specialization: "Sword",
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
		// Sense Skills
		"Direction": &Skill{
			Name: "Direction",
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
		"Eerie": &Skill{
			Name: "Eerie",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: sense,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Hearing": &Skill{
			Name: "Hearing",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: sense,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Sight": &Skill{
			Name: "Sight",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: sense,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Scrutinize": &Skill{
			Name: "Scrutinize",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: sense,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		// Knowledge Skills
		"Counterspell": &Skill{
			Name: "Counterspell",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: knowledge,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Healing": &Skill{
			Name: "Healing",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: knowledge,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Lore": &Skill{
			Name: "Lore",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: knowledge,
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
			LinkStat: knowledge,
			Dice: &DiePool{
				Normal: 0,
			},
			ReqSpec:        true,
			Specialization: "Elven",
		},
		"Strategy": &Skill{
			Name: "Strategy",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: knowledge,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Student of ": &Skill{
			Name: "Student of ",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: knowledge,
			Dice: &DiePool{
				Normal: 0,
			},
			ReqSpec:        true,
			Specialization: "Wyverns",
		},
		"Tactics": &Skill{
			Name: "Tactics",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: knowledge,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		// Charm Skills
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
		"Fascinate": &Skill{
			Name: "Fascinate",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: charm,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Graces": &Skill{
			Name: "Graces",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: charm,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Jest": &Skill{
			Name: "Jest",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: charm,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Plead": &Skill{
			Name: "Plead",
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
		"Haggle": &Skill{
			Name: "Haggle",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: command,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Inspire": &Skill{
			Name: "Inspire",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: command,
			Dice: &DiePool{
				Normal: 0,
			},
		},
		"Performing": &Skill{
			Name: "Performing",
			Quality: &Quality{
				Type:  "Useful",
				Level: 0,
			},
			LinkStat: command,
			Dice: &DiePool{
				Normal: 0,
			},
			ReqSpec:        true,
			Specialization: "Storyteller",
		},
		"Intimidate": &Skill{
			Name: "Intimidate",
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

	c.Advantages = []*Advantage{}

	return &c
}
