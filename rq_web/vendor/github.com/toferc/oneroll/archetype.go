package oneroll

import (
	"fmt"
)

// Archetype is a grouping of Sources, Permissions & Intrinsics that defines what powers a character can use
type Archetype struct {
	Type        string
	Sources     []*Source
	Permissions []*Permission
	Intrinsics  []*Intrinsic
	Cost        int
}

func (a Archetype) String() string {
	text := fmt.Sprintf("\nArchtype: %s (%dpts)", a.Type, a.Cost)

	text += "\nSources: "

	for _, s := range a.Sources {
		text += fmt.Sprintf("%s (%dpts), ", s.Type, s.Cost)
	}

	text = text[:len(text)-2]

	text += "\nPermissions: "

	for _, p := range a.Permissions {
		text += fmt.Sprintf("%s (%dpts), ", p.Type, p.Cost)
	}

	text = text[:len(text)-2]

	if len(a.Intrinsics) > 0 {
		text += "\nIntrinsics: "

		for _, i := range a.Intrinsics {
			text += fmt.Sprintf("%s (%dpts), ", i.Name, i.Cost)
		}
		text = text[:len(text)-2]
	}

	return text
}

// CalculateCost adds costs from sources, permissions and intrinsics
func (a *Archetype) CalculateCost() {
	var c int

	for _, s := range a.Sources {
		//First Source is free and all sources but focus cost 5pts
		if s.Type == a.Sources[0].Type && s.Cost > 0 {
			s.Cost = 0
		}
		c += s.Cost
	}

	for _, p := range a.Permissions {
		c += p.Cost
	}

	for _, i := range a.Intrinsics {
		if i.RequiresLevel {
			c += i.Cost * i.Level
		} else {
			c += i.Cost
		}
	}

	a.Cost = c
}

// Source is a source of a Character's powers
type Source struct {
	Type        string
	Cost        int // First source is free
	Description string
}

func (s Source) String() string {
	return fmt.Sprintf("%s (%dpts)", s.Type, s.Cost)
}

// Permission is the type of powers a Character can purchase
type Permission struct {
	Type              string
	Cost              int
	Description       string
	AllowHyperSkill   bool
	AllowHyperStat    bool
	ExceedStatLimit   bool
	AllowWiggle       bool
	AllowHard         bool
	AllowMiracles     bool
	AllowGadgeteering bool
	AllowGadgets      bool
	PowerLimit        int
}

func (p Permission) String() string {
	return fmt.Sprintf("%s (%dpts)", p.Type, p.Cost)
}

// Intrinsic is a modification from the human standard
type Intrinsic struct {
	Name          string
	Description   string
	RequiresLevel bool
	Level         int
	RequiresInfo  bool
	Info          string
	Cost          int
}

func (i Intrinsic) String() string {
	return fmt.Sprintf("%s (%dpts)", i.Name, i.Cost)
}

// Sources Set Wild Talents Sources
var Sources = map[string]Source{

	"Construct":        Source{Type: "Construct", Cost: 5, Description: ""},
	"Cyborg":           Source{Type: "Cyborg", Cost: 5, Description: ""},
	"Divine":           Source{Type: "Divine", Cost: 5, Description: ""},
	"Driven":           Source{Type: "Driven", Cost: 5, Description: ""},
	"Extraterrestrial": Source{Type: "Extraterrestrial", Cost: 5, Description: ""},
	"Genetic":          Source{Type: "Genetic", Cost: 5, Description: ""},
	"Life Force":       Source{Type: "Life Force", Cost: 5, Description: ""},
	"Paranormal":       Source{Type: "Paranormal", Cost: 5, Description: ""},
	"Power Focus":      Source{Type: "Power Focus", Cost: -8, Description: ""},
	"Psi":              Source{Type: "Psi", Cost: 5, Description: ""},
	"Technological":    Source{Type: "Technological", Cost: 5, Description: ""},
	"Unknown":          Source{Type: "Unknown", Cost: -5, Description: ""},
}

// Permissions sets Wild Talents default permissions
var Permissions = map[string]Permission{

	"None": Permission{
		Type:        "None",
		Cost:        0,
		Description: "",
	},
	"Hypertrained": Permission{
		Type:            "Hypertrained",
		Cost:            5,
		Description:     "",
		AllowHyperSkill: true,
	},
	"Inhuman Stats": Permission{
		Type:            "Inhuman Stats",
		Cost:            1,
		Description:     "",
		AllowHyperSkill: true,
	},
	"Inventor": Permission{
		Type:              "Inventor",
		Cost:              5,
		Description:       "",
		AllowGadgeteering: true,
		AllowGadgets:      true,
	},
	"One Power": Permission{
		Type:          "One Power",
		Cost:          1,
		Description:   "",
		AllowMiracles: true,
		PowerLimit:    1,
	},
	"Peak Performer": Permission{
		Type:        "Peak Performer",
		Cost:        5,
		Description: "",
		AllowWiggle: true,
		AllowHard:   true,
	},
	"Power Theme": Permission{
		Type:            "Power Theme",
		Cost:            5,
		Description:     "",
		AllowHyperSkill: true,
		AllowHyperStat:  true,
		AllowMiracles:   true,
		AllowHard:       true,
		AllowWiggle:     true,
	},
	"Prime Specimen": Permission{
		Type:           "Prime Specimen",
		Cost:           5,
		Description:    "",
		AllowHyperStat: true,
	},
	"Super": Permission{
		Type:            "Super",
		Cost:            15,
		Description:     "",
		AllowHyperSkill: true,
		AllowHyperStat:  true,
		AllowMiracles:   true,
		AllowHard:       true,
		AllowWiggle:     true,
	},
	"Super Equipment": Permission{
		Type:         "Super Equipment",
		Cost:         2,
		Description:  "",
		AllowGadgets: true,
	},
}

// Intrinsics sets Wild Talents default permissions
var Intrinsics = map[string]Intrinsic{

	"Allergy": Intrinsic{
		Name:          "Allergy",
		RequiresInfo:  true,
		Info:          "Substance frequency x threat",
		RequiresLevel: true,
		Level:         8,
		Description:   "Default to Frequent & Kills",
		Cost:          -1,
	},
	"Brute/Frail": Intrinsic{
		Name:        "Brute/Frail",
		Description: "",
		Cost:        -8,
	},
	"Custom Stats": Intrinsic{
		Name:         "Custom Stats",
		Description:  "",
		RequiresInfo: true,
		Info:         "",
		Cost:         5,
	},
	"Globular": Intrinsic{
		Name:        "Globular",
		Description: "",
		Cost:        8,
	},
	"Inhuman": Intrinsic{
		Name:          "Inhuman",
		RequiresInfo:  true,
		Info:          "",
		RequiresLevel: true,
		Level:         8,
		Description:   "Terrifying",
		Cost:          -1,
	},
	"Mandatory Power": Intrinsic{
		Name:         "Mandatory Power",
		RequiresInfo: true,
		Info:         "",
		Description:  "",
		Cost:         0,
	},
	"Mutable": Intrinsic{
		Name:        "Mutable",
		Description: "",
		Cost:        15,
	},
	"No Base Will": Intrinsic{
		Name:        "No Base Will",
		Description: "",
		Cost:        -10,
	},
	"No Willpower": Intrinsic{
		Name:        "No Willpower",
		Description: "",
		Cost:        -5,
	},
	"No Willpower No Way": Intrinsic{
		Name:        "No Willpower No Way",
		Description: "",
		Cost:        -5,
	},
	"Unhealing": Intrinsic{
		Name:        "Unhealing",
		Description: "",
		Cost:        -8,
	},
	"Vulnerable": Intrinsic{
		Name:          "Vulnerable",
		RequiresInfo:  true,
		Info:          "Substance frequency x threat",
		RequiresLevel: true,
		Level:         8,
		Description:   "Default to Frequent & Kills",
		Cost:          -1,
	},
	"Willpower Contest": Intrinsic{
		Name:        "Willpower Contest",
		Description: "",
		Cost:        -10,
	},
	"Custom": Intrinsic{
		Name:          "Custom",
		RequiresInfo:  true,
		Info:          "",
		RequiresLevel: true,
		Level:         1,
		Description:   "Default to Frequent & Kills",
		Cost:          -1,
	},
}
