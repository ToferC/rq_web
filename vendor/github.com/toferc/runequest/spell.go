package runequest

import "fmt"

// Spell is a magical effect
type Spell struct {
	Name        string
	Description string
	Domain      string
	CoreString  string
	UserString  string
	UserChoice  bool
	Range       int
	Duration    string
	Runes       []string
	Variable    bool
	Cost        int
	Points      int
	Source      string
}

// Domains represents domains of magic
var Domains = []string{"Rune", "Spirit", "Sorcery"}

func (s *Spell) String() string {

	text := fmt.Sprintf("%s (%dpts)", s.Name, s.Cost)

	return text
}

// GenerateName sets the skill map name
func (s *Spell) GenerateName() {

	var n string

	if s.UserString != "" {
		n = fmt.Sprintf("%s (%s)", s.CoreString, s.UserString)
	} else {
		n = s.CoreString
	}
	s.Name = n
}

func formatSpellArray(sa map[string]Spell) string {
	text := ""
	end := len(sa)
	counter := 0

	for _, v := range sa {
		if counter+1 == end {
			text += v.String()
		} else {
			text += v.String() + ", "
			counter++
		}
	}
	return text
}
