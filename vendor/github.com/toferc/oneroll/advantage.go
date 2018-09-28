package oneroll

import "fmt"

type Advantage struct {
	Name          string
	Description   string
	RequiresLevel bool
	Level         int
	RequiresInfo  bool
	Info          string
	Cost          int
}

func (a *Advantage) String() string {
	text := ""

	if a.RequiresLevel {
		text += fmt.Sprintf("\n%s (%d)",
			a.Name, a.Level*a.Cost)
	} else {
		text += fmt.Sprintf("%s (%d)",
			a.Name, a.Cost)
	}

	if a.RequiresInfo {
		text += fmt.Sprintf(" [%s]", a.Info)
	}

	return text
}

// Advantages sets Wild Talents default permissions
var Advantages = map[string]Advantage{
	"Animal Companion": Advantage{
		Name:          "Animal Companion",
		RequiresInfo:  true,
		Info:          "Type of Animal",
		RequiresLevel: true,
		Level:         1,
		Description:   "See p.29",
		Cost:          1,
	},
	"Beauty": Advantage{
		Name:          "Beauty",
		RequiresLevel: true,
		Level:         1,
		Description:   "See p.29",
		Cost:          1,
	},
	"Cannibal Smile": Advantage{
		Name:        "Cannibal Smile",
		Description: "See p.31",
		Cost:        1,
	},
	"Followers": Advantage{
		Name:          "Followers",
		RequiresInfo:  true,
		Info:          "Type of Follower",
		RequiresLevel: true,
		Level:         1,
		Description:   "See p.29",
		Cost:          1,
	},
	"Fool Lucky": Advantage{
		Name:        "Fool Lucky",
		Description: "See p.31",
		Cost:        5,
	},
	"Knack for Learning": Advantage{
		Name:         "Knack for Learning",
		Description:  "See p.31",
		RequiresInfo: true,
		Info:         "Skill",
		Cost:         5,
	},
	"Leather Hard": Advantage{
		Name:        "Leather Hard",
		Description: "See p.31",
		Cost:        5,
	},
	"Lucky": Advantage{
		Name:        "Lucky",
		Description: "See p.31",
		Cost:        1,
	},
	"Patron": Advantage{
		Name:          "Patron",
		RequiresInfo:  true,
		Info:          "Type of Patron",
		RequiresLevel: true,
		Level:         1,
		Description:   "See p.29",
		Cost:          1,
	},
	"Spells": Advantage{
		Name:          "Spells",
		RequiresInfo:  true,
		Info:          "Type of Spells",
		RequiresLevel: true,
		Level:         1,
		Description:   "See p.29",
		Cost:          1,
	},
	"Possession": Advantage{
		Name:          "Possession",
		RequiresInfo:  true,
		Info:          "Type of Possession",
		RequiresLevel: true,
		Level:         1,
		Description:   "See p.29",
		Cost:          1,
	},
	"Secret": Advantage{
		Name:          "Secret",
		RequiresInfo:  true,
		Info:          "Type of Secret",
		RequiresLevel: true,
		Level:         1,
		Description:   "See p.29",
		Cost:          1,
	},
	"Status": Advantage{
		Name:          "Status",
		RequiresInfo:  true,
		Info:          "Type of Status",
		RequiresLevel: true,
		Level:         1,
		Description:   "See p.29",
		Cost:          1,
	},
	"Thick Headed": Advantage{
		Name:        "Thick Headed",
		Description: "See p.33",
		Cost:        1,
	},
	"Wealth": Advantage{
		Name:          "Wealth",
		RequiresInfo:  true,
		Info:          "Type of Wealth",
		RequiresLevel: true,
		Level:         1,
		Description:   "See p.29",
		Cost:          1,
	},
}
