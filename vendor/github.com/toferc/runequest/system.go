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

// IDCoreRunes returns the top 3 runes for the character
func (c *Character) IDCoreRunes() {

	// Reset CoreRunes
	c.CoreRunes = []*Ability{}

	var runes []*Ability

	for _, e := range c.ElementalRunes {
		e.UpdateAbility()
		if e.Total > 50 {
			runes = append(runes, e)
		}
	}

	for _, p := range c.PowerRunes {
		p.UpdateAbility()
		if p.Total > 50 {
			runes = append(runes, p)
		}
	}

	for _, c := range c.ConditionRunes {
		c.UpdateAbility()
		if c.Total > 50 {
			runes = append(runes, c)
		}
	}

	// Sort Runes
	sort.Sort(ByTotal(runes))

	// Return max of 3 runes
	l := len(runes)

	if l > 3 {
		l = 3
	}

	for _, r := range runes[:l] {
		c.CoreRunes = append(c.CoreRunes, r)
	}
}

// DetermineRuneModifiers adds stat modifiers based on runes
func (c *Character) DetermineRuneModifiers() []string {

	var runes []*Ability

	var runeModifiers []string

	// Add abilities to array for sorting
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
		runeModifiers = append(runeModifiers, "STR")
	case primary == "Earth":
		runeModifiers = append(runeModifiers, "CON")
	case primary == "Darkness":
		runeModifiers = append(runeModifiers, "SIZ")
	case primary == "Fire/Sky":
		runeModifiers = append(runeModifiers, "INT")
	case primary == "Water":
		runeModifiers = append(runeModifiers, "DEX")
	case primary == "Moon":
		runeModifiers = append(runeModifiers, "POW")
	}

	switch {
	case secondary == "Air":
		runeModifiers = append(runeModifiers, "STR")
	case secondary == "Earth":
		runeModifiers = append(runeModifiers, "CON")
	case secondary == "Darkness":
		runeModifiers = append(runeModifiers, "SIZ")
	case secondary == "Fire/Sky":
		runeModifiers = append(runeModifiers, "INT")
	case secondary == "Water":
		runeModifiers = append(runeModifiers, "DEX")
	case secondary == "Moon":
		runeModifiers = append(runeModifiers, "POW")
	}

	fmt.Println(runeModifiers)
	return runeModifiers
}
