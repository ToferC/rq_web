package oneroll

import (
	"fmt"
)

// Statistic represents common attributes possessed by every character
type Statistic struct {
	Name        string
	Dice        *DiePool
	HyperStat   *HyperStat
	Modifiers   []*Modifier
	EffectsWill bool
	Cost        int
}

// HyperStat is a modified version of a regular Statistic
type HyperStat struct {
	Name       string
	Qualities  []*Quality
	Dice       *DiePool
	Effect     string
	Apply      bool
	CostPerDie int
	Cost       int
}

func (s *Statistic) getDiePool() *DiePool {

	td := &DiePool{}

	if s.HyperStat != nil {

		td.Normal = s.Dice.Normal + s.HyperStat.Dice.Normal
		td.Hard = s.Dice.Hard + s.HyperStat.Dice.Hard
		td.Wiggle = s.Dice.Wiggle + s.HyperStat.Dice.Wiggle

		for _, q := range s.HyperStat.Qualities {
			for _, m := range q.Modifiers {
				if m.Name == "Spray" {
					td.Spray = m.Level
				}

				if m.Name == "Go First" {
					td.GoFirst = m.Level
				}
			}
		}
	} else {
		td = s.Dice
	}
	return td
}

// FormatDiePool returns a die string
func (s *Statistic) FormatDiePool(actions int) string {

	td := ReturnDice(s)

	normal := td.Normal
	hard := td.Hard
	wiggle := td.Wiggle
	goFirst := td.GoFirst
	spray := td.Spray

	text := fmt.Sprintf("%dac+%dd+%dhd+%dwd+%dgf+%dsp",
		actions,
		normal,
		hard,
		wiggle,
		goFirst,
		spray)

	return text
}

func (s Statistic) String() string {

	td := ReturnDice(&s)

	text := ""

	if s.HyperStat != nil {
		text += fmt.Sprintf("%s*: %s", s.Name, td)
	} else {
		text = fmt.Sprintf("%s: %s", s.Name, td)
	}

	return text
}

func (hs HyperStat) String() string {
	text := fmt.Sprintf("%s %s (", hs.Name, hs.Dice)

	for _, q := range hs.Qualities {
		text += fmt.Sprintf("%s", string(q.Type[0]))
		if q.Level > 0 {
			text += fmt.Sprintf("+%d", q.Level)
		}
	}

	text += fmt.Sprintf(") [%d/die] %dpts\n",
		hs.CostPerDie,
		hs.Cost)

	return text
}

// CalculateCost determines the cost of a Power Quality
// Called from Character.CalculateCharacterCost()
func (s *Statistic) CalculateCost() {
	// Base Cost
	b := 5
	// Modifier Cost
	mc := 0

	// Add costs for modifiers
	for _, m := range s.Modifiers {
		m.CalculateCost(0)
		if m.RequiresLevel {
			mc += m.CostPerLevel * m.Level
		} else {
			mc += m.CostPerLevel
		}
	}

	// Add cost for HyperStat levels applied to Stat
	if s.HyperStat != nil {
		for _, q := range s.HyperStat.Qualities {
			if q.Level > 0 {
				b += q.Level
			}
		}
	}

	if len(s.Modifiers) > 0 && mc < 1 {
		// There are mods, but flaws reduce cost below 1
		mc = 1
	}

	b += mc

	total := b * s.Dice.Normal
	total += b * 2 * s.Dice.Hard
	total += b * 4 * s.Dice.Wiggle

	s.Cost = total
}

// CalculateCost generates and udpates the cost for HypeSKills
func (hs *HyperStat) CalculateCost() {

	b := 4 // Base Cost

	for _, q := range hs.Qualities {

		// Add Power Capacity Modifier if needed
		if len(q.Capacities) > 3 {
			tm := Modifiers["Power Capacity"]
			tm.Level = len(q.Capacities) - 3
			q.Modifiers = append(q.Modifiers, &tm)
		}

		for _, m := range q.Modifiers {
			m.CalculateCost(0)
		}
		q.CalculateCost(0)
		b += q.CostPerDie
	}

	hs.CostPerDie = b

	total := b * hs.Dice.Normal
	total += b * 2 * hs.Dice.Hard
	total += b * 4 * hs.Dice.Wiggle

	hs.Cost = total
}
