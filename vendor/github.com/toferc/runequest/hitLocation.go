package runequest

import "fmt"

// HitLocation represents a body area that can take damage
type HitLocation struct {
	Name     string
	HitLoc   []int
	Base     int
	Min      int
	Max      int
	Value    int
	Updates  []*Update
	Wounds   []bool
	Armor    int
	Disabled bool
	Maimed   bool
}

// GenerateHitLocationMap takes a HitLocation map and generates an array of strings
func GenerateHitLocationMap(hlForm map[string]*HitLocation) []string {
	m := []string{}
	for k := range hlForm {
		m = append(m, k)
	}
	return m
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

// LocationForms is a map of standard hit locations
var LocationForms = map[string]map[string]*HitLocation{
	"Humanoids":           HumanoidLocations,
	"Centaurs":            CentaurLocations,
	"Dragons/Manticores":  DragonHitLocations,
	"Four-Legged Animals": FourLeggedAnimalsHitLocations,
}

// HumanoidLocations is a base map of hit locations
var HumanoidLocations = map[string]*HitLocation{
	"R Leg": &HitLocation{
		Name:   "R Leg",
		HitLoc: []int{1, 2, 3, 4},
		Base:   0,
	},
	"L Leg": &HitLocation{
		Name:   "L Leg",
		HitLoc: []int{5, 6, 7, 8},
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

// FourLeggedAnimalsHitLocations is a base map of hit locations
var FourLeggedAnimalsHitLocations = map[string]*HitLocation{
	"L Hind Leg": &HitLocation{
		Name:   "L Hind Leg",
		HitLoc: []int{3, 4},
		Base:   0,
	},
	"R Hind Leg": &HitLocation{
		Name:   "R Hind Leg",
		HitLoc: []int{1, 2},
		Base:   0,
	},
	"Hindquarters": &HitLocation{
		Name:   "Hindquarters",
		HitLoc: []int{5, 6, 7},
		Base:   1,
	},
	"Forequarters": &HitLocation{
		Name:   "Forequarters",
		HitLoc: []int{8, 9, 10},
		Base:   1,
	},
	"L Foreleg": &HitLocation{
		Name:   "L Foreleg",
		HitLoc: []int{11, 12, 13},
		Base:   0,
	},
	"R Foreleg": &HitLocation{
		Name:   "R Foreleg",
		HitLoc: []int{14, 15, 16},
		Base:   0,
	},
	"Head": &HitLocation{
		Name:   "Head",
		HitLoc: []int{17, 18, 19, 20},
		Base:   0,
	},
}

// DragonHitLocations is a base map of hit locations
var DragonHitLocations = map[string]*HitLocation{
	"Tail": &HitLocation{
		Name:   "Tail",
		HitLoc: []int{1, 2},
		Base:   0,
	},
	"L Hind Leg": &HitLocation{
		Name:   "L Hind Leg",
		HitLoc: []int{5, 6},
		Base:   0,
	},
	"R Hind Leg": &HitLocation{
		Name:   "R Hind Leg",
		HitLoc: []int{3, 4},
		Base:   0,
	},
	"Hindquarters": &HitLocation{
		Name:   "Hindquarters",
		HitLoc: []int{7, 8},
		Base:   1,
	},
	"Forequarters": &HitLocation{
		Name:   "Forequarters",
		HitLoc: []int{9, 10},
		Base:   1,
	},
	"R Wing": &HitLocation{
		Name:   "R Wing",
		HitLoc: []int{11, 12},
		Base:   -1,
	},
	"L Wing": &HitLocation{
		Name:   "L Wing",
		HitLoc: []int{13, 14},
		Base:   -1,
	},
	"L Foreleg": &HitLocation{
		Name:   "L Foreleg",
		HitLoc: []int{15, 16},
		Base:   0,
	},
	"R Foreleg": &HitLocation{
		Name:   "R Foreleg",
		HitLoc: []int{17, 18},
		Base:   0,
	},
	"Head": &HitLocation{
		Name:   "Head",
		HitLoc: []int{19, 20},
		Base:   0,
	},
}

// CentaurLocations is a base map of hit locations
var CentaurLocations = map[string]*HitLocation{
	"R Hind Leg": &HitLocation{
		Name:   "R Hind Leg",
		HitLoc: []int{1, 2},
		Base:   -1,
	},
	"L Hind Leg": &HitLocation{
		Name:   "L Hind Leg",
		HitLoc: []int{3, 4},
		Base:   -1,
	},
	"Hindquarters": &HitLocation{
		Name:   "Hindquarters",
		HitLoc: []int{5, 6},
		Base:   1,
	},
	"Forequarters": &HitLocation{
		Name:   "Forequarters",
		HitLoc: []int{7, 8},
		Base:   1,
	},
	"R Foreleg": &HitLocation{
		Name:   "R Foreleg",
		HitLoc: []int{9, 10},
		Base:   -1,
	},
	"L Foreleg": &HitLocation{
		Name:   "L Foreleg",
		HitLoc: []int{11, 12},
		Base:   -1,
	},
	"Chest": &HitLocation{
		Name:   "Chest",
		HitLoc: []int{13, 14},
		Base:   1,
	},
	"R Arm": &HitLocation{
		Name:   "R Arm",
		HitLoc: []int{15, 16},
		Base:   -1,
	},
	"L Arm": &HitLocation{
		Name:   "L Arm",
		HitLoc: []int{17, 18},
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
