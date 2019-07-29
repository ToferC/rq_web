package runequest

import "fmt"

// Movement represents a movement type and value
type Movement struct {
	Name  string
	Value int
}

func (m *Movement) String() string {
	text := fmt.Sprintf("%s: %d", m.Name, m.Value)
	return text
}
