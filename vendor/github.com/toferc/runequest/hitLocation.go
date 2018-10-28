package runequest

import "fmt"

// HitLocation represents a body area that can take damage
type HitLocation struct {
	Name     string
	HitLoc   []int
	Max      int
	Value    int
	Updates  []*Update
	Wounds   []bool
	Armor    int
	Disabled bool
}

// Strings
func (l HitLocation) String() string {
	text := fmt.Sprintf("(%s) - %s\n",
		TrimSliceBrackets(l.HitLoc),
		l.Name,
	)

	if l.Armor > 0 {
		text += fmt.Sprintf("Armor: %d ", l.Armor)
	}

	v := l.CountWounds()

	if v > 0 {
		text += fmt.Sprintf(" HP: %d", v)
	}

	return text
}

// FillWounds creates the array of empty wound boxes
func (l *HitLocation) FillWounds() {
	for i := 0; i < l.Value; i++ {
		l.Wounds = append(l.Wounds, false)
	}
}

// CountWounds returns the number of wounds based on the bool Kill & Shock slices
func (l *HitLocation) CountWounds() int {

	var hp int

	for _, box := range l.Wounds {
		if box {
			hp++
		}
	}

	return hp
}
