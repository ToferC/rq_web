package oneroll

import "fmt"

// Location represents a body area that can take damage
type Location struct {
	Name     string
	HitLoc   []int
	Boxes    int
	Shock    []bool
	Kill     []bool
	LAR      int
	HAR      int
	Disabled bool
}

// Strings
func (l Location) String() string {
	text := fmt.Sprintf("(%s) - %s\n",
		TrimSliceBrackets(l.HitLoc),
		l.Name,
	)

	if l.LAR > 0 {
		text += fmt.Sprintf("LAR: %d ", l.LAR)
	}

	if l.HAR > 0 {
		text += fmt.Sprintf("HAR: %d ", l.HAR)
	}

	k, s := l.CountWounds()

	if k > 0 {
		text += fmt.Sprintf(" Kill: %d", k)
	}

	if s > 0 {
		text += fmt.Sprintf(" Shock: %d", s)
	}
	return text
}

// FillWounds creates the array of empty wound boxes
func (l *Location) FillWounds() {
	for i := 0; i < l.Boxes; i++ {
		l.Kill = append(l.Kill, false)
		l.Shock = append(l.Shock, false)
	}
}

// CountWounds returns the number of wounds based on the bool Kill & Shock slices
func (l *Location) CountWounds() (int, int) {

	var kill, shock int

	for _, box := range l.Kill {
		if box {
			kill++
		}
	}

	for _, box := range l.Shock {
		if box {
			shock++
		}
	}
	return kill, shock
}
