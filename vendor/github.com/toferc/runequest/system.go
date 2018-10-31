package runequest

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// RollDice rolls and sum dice
func RollDice(max, min, bonus, numDice int) int {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	result := 0
	for i := 1; i < numDice+1; i++ {
		roll := r1.Intn(max+1-min) + min
		result += roll
	}

	result += bonus

	return result
}

// Sorting Functions

// ByTotal implements the sort interface for abilities
type ByTotal []*Ability

func (a ByTotal) Len() int           { return len(a) }
func (a ByTotal) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTotal) Less(i, j int) bool { return a[i].Total > a[j].Total }

// AddRuneModifiers adds stat modifiers based on runes
func (c *Character) AddRuneModifiers() {

	var runes []*Ability

	for _, a := range c.ElementalRunes {
		a.UpdateAbility()
		runes = append(runes, a)
	}

	// Reset Rune Bonuses
	for _, v := range c.Statistics {
		v.RuneBonus = 0
	}

	// Sort Runes
	sort.Sort(ByTotal(runes))
	fmt.Println(runes)

	primary, secondary := runes[0].Name, runes[1].Name

	switch {
	case primary == "Air":
		c.Statistics["STR"].RuneBonus = 2
	case primary == "Earth":
		c.Statistics["CON"].RuneBonus = 2
	case primary == "Darkness":
		c.Statistics["SIZ"].RuneBonus = 2
	case primary == "Fire/Sky":
		c.Statistics["INT"].RuneBonus = 2
	case primary == "Water":
		c.Statistics["DEX"].RuneBonus = 2
	case primary == "Moon":
		c.Statistics["POW"].RuneBonus = 2
	}

	switch {
	case secondary == "Air":
		c.Statistics["STR"].RuneBonus = 1
	case secondary == "Earth":
		c.Statistics["CON"].RuneBonus = 1
	case secondary == "Darkness":
		c.Statistics["SIZ"].RuneBonus = 1
	case secondary == "Fire/Sky":
		c.Statistics["INT"].RuneBonus = 1
	case secondary == "Water":
		c.Statistics["DEX"].RuneBonus = 1
	case secondary == "Moon":
		c.Statistics["POW"].RuneBonus = 1
	}

}
