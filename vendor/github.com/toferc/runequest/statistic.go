package runequest

import "fmt"

// Statistic is a core element for all characters in an RPG system
type Statistic struct {
	Name            string
	Value           int
	Base            int
	RuneBonus       int
	Updates         []*Update
	Total           int
	Max             int
	Min             int
	ExperienceCheck bool
}

func (s *Statistic) String() string {

	s.Total = s.Base + s.Value + s.RuneBonus

	text := fmt.Sprintf("%s: %d", s.Name, s.Total)
	return text
}

// StatMap gives default ordering for stats
var StatMap = []string{
	"STR", "DEX", "CON", "SIZ", "POW", "INT", "CHA"}

// TotalStatistics updates values for stats after being modified
func (c *Character) TotalStatistics() {

	for _, s := range c.Statistics {

		updates := 0

		for _, u := range s.Updates {
			updates += u.Value
		}

		s.Total = s.Base + s.Value + s.RuneBonus + updates
	}
}

// RuneQuestStats is the base stats for RuneQuest
var RuneQuestStats = map[string]*Statistic{
	"STR": &Statistic{
		Name: "Strength",
		Base: RollDice(6, 1, 0, 3),
	},
	"DEX": &Statistic{
		Name: "Dexterity",
		Base: RollDice(6, 1, 0, 3),
	},
	"INT": &Statistic{
		Name: "Intelligence",
		Base: RollDice(6, 1, 6, 2),
	},
	"CON": &Statistic{
		Name: "Constitution",
		Base: RollDice(6, 1, 6, 2),
	},
	"POW": &Statistic{
		Name: "Power",
		Base: RollDice(6, 1, 0, 3),
	},
	"SIZ": &Statistic{
		Name: "Size",
		Base: RollDice(6, 1, 6, 2),
	},
	"CHA": &Statistic{
		Name: "Charisma",
		Base: RollDice(6, 1, 0, 3),
	},
}
