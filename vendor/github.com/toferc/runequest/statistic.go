package runequest

import "fmt"

// Statistic is a core element for all characters in an RPG system
type Statistic struct {
	Name            string
	Base            int
	RuneBonus       int
	HomelandBonus   int
	Updates         []*Update
	Total           int
	Max             int
	Min             int
	ExperienceCheck bool
}

func (s *Statistic) String() string {

	text := fmt.Sprintf("%s: %d", s.Name, s.Total)
	return text
}

// StatMap gives default ordering for stats
var StatMap = []string{
	"STR", "CON", "SIZ", "DEX", "INT", "POW", "CHA"}

// SpiritMap is a base map of stats for Spirits
var SpiritMap = []string{
	"INT", "POW", "CHA",
}

// ElementalMap is a base map of stast for Elementals
var ElementalMap = []string{
	"SIZ", "POW", "STR",
}

// UpdateStatistic updates values for stats after being modified
func (s *Statistic) UpdateStatistic() {

	updates := 0

	for _, u := range s.Updates {
		updates += u.Value
	}

	s.Total = s.Base + s.RuneBonus + updates

	if s.Total > s.Max {
		s.Total = s.Max
	}
}

// RuneQuestStats is the base stats for RuneQuest
var RuneQuestStats = map[string]*Statistic{
	"STR": &Statistic{
		Name: "Strength",
		Base: RollDice(6, 1, 0, 3),
		Max:  18,
	},
	"DEX": &Statistic{
		Name: "Dexterity",
		Base: RollDice(6, 1, 0, 3),
		Max:  18,
	},
	"INT": &Statistic{
		Name: "Intelligence",
		Base: RollDice(6, 1, 6, 2),
		Max:  18,
	},
	"CON": &Statistic{
		Name: "Constitution",
		Base: RollDice(6, 1, 6, 2),
		Max:  18,
	},
	"POW": &Statistic{
		Name: "Power",
		Base: RollDice(6, 1, 0, 3),
		Max:  18,
	},
	"SIZ": &Statistic{
		Name: "Size",
		Base: RollDice(6, 1, 6, 2),
		Max:  18,
	},
	"CHA": &Statistic{
		Name: "Charisma",
		Base: RollDice(6, 1, 0, 3),
		Max:  18,
	},
}
