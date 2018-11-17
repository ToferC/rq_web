package runequest

import "fmt"

// HitLocation represents a body area that can take damage
type HitLocation struct {
	Name     string
	HitLoc   []int
	Base     int
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

	text += fmt.Sprintf(" HP: %d/%d", l.Value, l.Max)

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

// Locations is a base map of hit locations
var Locations = map[string]*HitLocation{
	"L Leg": &HitLocation{
		Name:   "L Leg",
		HitLoc: []int{5, 6, 7, 8},
		Base:   0,
	},
	"R Leg": &HitLocation{
		Name:   "R Leg",
		HitLoc: []int{1, 2, 3, 4},
		Base:   0,
	},
	"Abdomen": &HitLocation{
		Name:   "Abdomen",
		HitLoc: []int{9, 10, 11},
		Base:   0,
	},
	"Chest": &HitLocation{
		Name:   "Chest",
		HitLoc: []int{12},
		Base:   1,
	},
	"L Arm": &HitLocation{
		Name:   "L Arm",
		HitLoc: []int{16, 17, 18},
		Base:   -1,
	},
	"R Arm": &HitLocation{
		Name:   "R Arm",
		HitLoc: []int{13, 14, 15},
		Base:   -1,
	},
	"Head": &HitLocation{
		Name:   "Head",
		HitLoc: []int{19, 20},
		Base:   0,
	},
}

// HPPerLocation is the base wound map
var HPPerLocation = []int{2, 3, 4, 5, 6, 7}
