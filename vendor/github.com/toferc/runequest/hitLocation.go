package runequest

import (
	"fmt"
	"sort"
)

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
	"Humanoids":                   HumanoidLocations,
	"Humanoids, Winged":           HumanoidWingedLocations,
	"Centaurs":                    CentaurLocations,
	"Dragons/Manticores":          DragonHitLocations,
	"Dragonewts":                  DragonewtHitLocations,
	"Four-Legged Animals":         FourLeggedAnimalsHitLocations,
	"Four-Legged Animals, Winged": FourLeggedAnimalsWingedHitLocations,
	"Serpents":                    SerpentLocations,
	"Serpents, Winged":            SerpentWingedLocations,
	"Birds, Flying":               BirdsFlyingHitLocations,
	"Birds, Running":              BirdsRunningHitLocations,
	"Beetles":                     BeetlesHitLocations,
	"Wyverns":                     WyvernsHitLocations,
	"Spiders, Giant":              SpidersGiantHitLocations,
	"Mammoths/Mastodons":          MammothsHitLocations,
	"Spirits":                     SpiritHitLocations,
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

// SerpentLocations is a base map of hit locations
var SerpentLocations = map[string]*HitLocation{
	"Tail": &HitLocation{
		Name:   "Tail",
		HitLoc: []int{1, 2, 3, 4, 5, 6},
		Base:   0,
	},
	"Body": &HitLocation{
		Name:   "Body",
		HitLoc: []int{7, 8, 9, 10, 11, 12, 13, 14},
		Base:   1,
	},
	"Head": &HitLocation{
		Name:   "Head",
		HitLoc: []int{15, 16, 17, 18, 19, 20},
		Base:   0,
	},
}

// SerpentWingedLocations is a base map of hit locations
var SerpentWingedLocations = map[string]*HitLocation{
	"Tail": &HitLocation{
		Name:   "Tail",
		HitLoc: []int{1, 2, 3, 4},
		Base:   0,
	},
	"Abdomen": &HitLocation{
		Name:   "Abdomen",
		HitLoc: []int{5, 6, 7, 8},
		Base:   0,
	},
	"Chest": &HitLocation{
		Name:   "Chest",
		HitLoc: []int{9, 10, 11, 12},
		Base:   1,
	},
	"R Wing": &HitLocation{
		Name:   "R Wing",
		HitLoc: []int{13, 14},
		Base:   -1,
	},
	"L Wing": &HitLocation{
		Name:   "L Wing",
		HitLoc: []int{15, 16},
		Base:   -1,
	},
	"Head": &HitLocation{
		Name:   "Head",
		HitLoc: []int{17, 18, 19, 20},
		Base:   0,
	},
}

// HumanoidWingedLocations is a base map of hit locations
var HumanoidWingedLocations = map[string]*HitLocation{
	"R Leg": &HitLocation{
		Name:   "R Leg",
		HitLoc: []int{1, 2, 3},
		Base:   0,
	},
	"L Leg": &HitLocation{
		Name:   "L Leg",
		HitLoc: []int{4, 5, 6},
		Base:   0,
	},
	"Abdomen": &HitLocation{
		Name:   "Abdomen",
		HitLoc: []int{7, 8, 9},
		Base:   0,
	},
	"Chest": &HitLocation{
		Name:   "Chest",
		HitLoc: []int{10},
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

// FourLeggedAnimalsHitLocations is a base map of hit locations
var FourLeggedAnimalsHitLocations = map[string]*HitLocation{
	"R Hind Leg": &HitLocation{
		Name:   "R Hind Leg",
		HitLoc: []int{1, 2},
		Base:   0,
	},
	"L Hind Leg": &HitLocation{
		Name:   "L Hind Leg",
		HitLoc: []int{3, 4},
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

// FourLeggedAnimalsWingedHitLocations is a base map of hit locations
var FourLeggedAnimalsWingedHitLocations = map[string]*HitLocation{
	"R Hind Leg": &HitLocation{
		Name:   "R Hind Leg",
		HitLoc: []int{1, 2, 3},
		Base:   0,
	},
	"L Hind Leg": &HitLocation{
		Name:   "L Hind Leg",
		HitLoc: []int{4, 5, 6},
		Base:   0,
	},
	"Hindquarters": &HitLocation{
		Name:   "Hindquarters",
		HitLoc: []int{7, 8, 9},
		Base:   1,
	},
	"Forequarters": &HitLocation{
		Name:   "Forequarters",
		HitLoc: []int{10},
		Base:   1,
	},
	"L Wing": &HitLocation{
		Name:   "L Wing",
		HitLoc: []int{11, 12},
		Base:   -1,
	},
	"R Wing": &HitLocation{
		Name:   "R Wing",
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

// MammothsHitLocations is a base map of hit locations
var MammothsHitLocations = map[string]*HitLocation{
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
		Base:   -1,
	},
	"R Foreleg": &HitLocation{
		Name:   "R Foreleg",
		HitLoc: []int{14, 15, 16},
		Base:   -1,
	},
	"Trunk": &HitLocation{
		Name:   "Trunk",
		HitLoc: []int{17},
		Base:   -3,
	},
	"Head": &HitLocation{
		Name:   "Head",
		HitLoc: []int{18, 19, 20},
		Base:   1,
	},
}

// BeetlesHitLocations is a base map of hit locations
var BeetlesHitLocations = map[string]*HitLocation{
	"R Hind Leg": &HitLocation{
		Name:   "R Hind Leg",
		HitLoc: []int{1},
		Base:   -1,
	},
	"R Center Leg": &HitLocation{
		Name:   "R Center Leg",
		HitLoc: []int{2},
		Base:   -1,
	},
	"L Hind Leg": &HitLocation{
		Name:   "L Hind Leg",
		HitLoc: []int{3},
		Base:   -1,
	},
	"L Center Leg": &HitLocation{
		Name:   "L Center Leg",
		HitLoc: []int{4},
		Base:   -1,
	},
	"Abdomen": &HitLocation{
		Name:   "Abdomen",
		HitLoc: []int{5, 6, 7, 8},
		Base:   3,
	},
	"Thorax": &HitLocation{
		Name:   "Thorax",
		HitLoc: []int{9, 10, 11, 12},
		Base:   3,
	},
	"R Foreleg": &HitLocation{
		Name:   "R Foreleg",
		HitLoc: []int{13, 14},
		Base:   -1,
	},
	"L Foreleg": &HitLocation{
		Name:   "L Foreleg",
		HitLoc: []int{15, 16},
		Base:   -1,
	},
	"Head": &HitLocation{
		Name:   "Head",
		HitLoc: []int{17, 18, 19, 20},
		Base:   2,
	},
}

// SpidersGiantHitLocations is a base map of hit locations
var SpidersGiantHitLocations = map[string]*HitLocation{
	"R Hind Leg": &HitLocation{
		Name:   "R Hind Leg",
		HitLoc: []int{1},
		Base:   -2,
	},
	"L Hind Leg": &HitLocation{
		Name:   "L Hind Leg",
		HitLoc: []int{2},
		Base:   -2,
	},
	"R Center Rear Leg": &HitLocation{
		Name:   "R Center Rear Leg",
		HitLoc: []int{3},
		Base:   -2,
	},
	"L Center Rear Leg": &HitLocation{
		Name:   "L Center Rear Leg",
		HitLoc: []int{4},
		Base:   -2,
	},
	"Abdomen": &HitLocation{
		Name:   "Abdomen",
		HitLoc: []int{5, 6, 7, 8},
		Base:   1,
	},
	"R Center Front Leg": &HitLocation{
		Name:   "R Center Front Leg",
		HitLoc: []int{9, 10},
		Base:   -2,
	},
	"L Center Front Leg": &HitLocation{
		Name:   "L Center Front Leg",
		HitLoc: []int{11, 12},
		Base:   -2,
	},
	"R Foreleg": &HitLocation{
		Name:   "R Foreleg",
		HitLoc: []int{13, 14},
		Base:   -2,
	},
	"L Foreleg": &HitLocation{
		Name:   "L Foreleg",
		HitLoc: []int{15, 16},
		Base:   -2,
	},
	"Head": &HitLocation{
		Name:   "Head",
		HitLoc: []int{17, 18, 19, 20},
		Base:   1,
	},
}

// WyvernsHitLocations is a base map of hit locations
var WyvernsHitLocations = map[string]*HitLocation{
	"Tail": &HitLocation{
		Name:   "Tail",
		HitLoc: []int{1},
		Base:   0,
	},
	"R Leg": &HitLocation{
		Name:   "R Leg",
		HitLoc: []int{2, 3, 4},
		Base:   0,
	},
	"L Leg": &HitLocation{
		Name:   "L Leg",
		HitLoc: []int{5, 6, 7},
		Base:   0,
	},
	"Abdomen": &HitLocation{
		Name:   "Abdomen",
		HitLoc: []int{8, 9},
		Base:   0,
	},
	"Chest": &HitLocation{
		Name:   "Chest",
		HitLoc: []int{10, 11, 12},
		Base:   1,
	},
	"R Wing": &HitLocation{
		Name:   "R Wing",
		HitLoc: []int{13, 14},
		Base:   -1,
	},
	"L Wing": &HitLocation{
		Name:   "L Wing",
		HitLoc: []int{15, 16},
		Base:   -1,
	},
	"Head": &HitLocation{
		Name:   "Head",
		HitLoc: []int{17, 18, 19, 20},
		Base:   0,
	},
}

// BirdsFlyingHitLocations is a base map of hit locations
var BirdsFlyingHitLocations = map[string]*HitLocation{
	"Tail": &HitLocation{
		Name:   "Tail",
		HitLoc: []int{1},
		Base:   -2,
	},
	"R Hind Leg": &HitLocation{
		Name:   "R Hind Leg",
		HitLoc: []int{2, 3, 4},
		Base:   -1,
	},
	"L Hind Leg": &HitLocation{
		Name:   "L Hind Leg",
		HitLoc: []int{5, 6, 7},
		Base:   -1,
	},
	"Abdomen": &HitLocation{
		Name:   "Abdomen",
		HitLoc: []int{8, 9},
		Base:   0,
	},
	"Chest": &HitLocation{
		Name:   "Chest",
		HitLoc: []int{10, 11, 12},
		Base:   1,
	},
	"R Wing": &HitLocation{
		Name:   "R Wing",
		HitLoc: []int{13, 14},
		Base:   0,
	},
	"L Wing": &HitLocation{
		Name:   "L Wing",
		HitLoc: []int{15, 16},
		Base:   0,
	},
	"Head": &HitLocation{
		Name:   "Head",
		HitLoc: []int{17, 18, 19, 20},
		Base:   0,
	},
}

// BirdsRunningHitLocations is a base map of hit locations
var BirdsRunningHitLocations = map[string]*HitLocation{
	"R Hind Leg": &HitLocation{
		Name:   "R Hind Leg",
		HitLoc: []int{1, 2, 3, 4},
		Base:   0,
	},
	"L Hind Leg": &HitLocation{
		Name:   "L Hind Leg",
		HitLoc: []int{5, 6, 7, 8},
		Base:   0,
	},
	"Abdomen": &HitLocation{
		Name:   "Abdomen",
		HitLoc: []int{9, 10},
		Base:   0,
	},
	"Chest": &HitLocation{
		Name:   "Chest",
		HitLoc: []int{11, 12, 13},
		Base:   1,
	},
	"R Wing": &HitLocation{
		Name:   "R Wing",
		HitLoc: []int{14, 15},
		Base:   -1,
	},
	"L Wing": &HitLocation{
		Name:   "L Wing",
		HitLoc: []int{16, 17},
		Base:   -1,
	},
	"Head": &HitLocation{
		Name:   "Head",
		HitLoc: []int{18, 19, 20},
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

// DragonewtHitLocations is a base map of hit locations
var DragonewtHitLocations = map[string]*HitLocation{
	"Tail": &HitLocation{
		Name:   "Tail",
		HitLoc: []int{1, 2},
		Base:   0,
	},
	"R Leg": &HitLocation{
		Name:   "R Leg",
		HitLoc: []int{3, 4, 5},
		Base:   0,
	},
	"L Leg": &HitLocation{
		Name:   "L Leg",
		HitLoc: []int{6, 7, 8},
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

// SpiritHitLocations is a base map of hit locations
var SpiritHitLocations = map[string]*HitLocation{
	"Spirit": &HitLocation{
		Name:   "Spirit",
		HitLoc: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		Base:   0,
	},
}

// SortLocations HitLocations
func SortLocations(locations map[string]*HitLocation) []string {
	locationArray := []*HitLocation{}

	for _, v := range locations {
		locationArray = append(locationArray, v)
	}

	hitloc := func(hl1, hl2 *HitLocation) bool {
		return hl1.HitLoc[0] > hl2.HitLoc[0]
	}

	ByHL(hitloc).SortHL(locationArray)

	stringArray := []string{}

	for _, l := range locationArray {
		stringArray = append(stringArray, l.Name)
	}

	return stringArray
}

// ByHL is the type of a "less" function that defines the ordering of its HitLoc[0] arguments.
type ByHL func(hl1, hl2 *HitLocation) bool

// SortHL is a method on the function type, By, that sorts the argument slice according to the function.
func (by ByHL) SortHL(locations []*HitLocation) {
	ls := &locationsorter{
		locations: locations,
		by:        by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ls)
}

// locationsorter joins a By function and a slice of Planets to be sorted.
type locationsorter struct {
	locations []*HitLocation
	by        func(hl1, hl2 *HitLocation) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (hl *locationsorter) Len() int {
	return len(hl.locations)
}

// Swap is part of sort.Interface.
func (hl *locationsorter) Swap(i, j int) {
	hl.locations[i], hl.locations[j] = hl.locations[j], hl.locations[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (hl *locationsorter) Less(i, j int) bool {
	return hl.by(hl.locations[i], hl.locations[j])
}
