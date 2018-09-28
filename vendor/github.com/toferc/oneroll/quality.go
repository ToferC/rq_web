package oneroll

import (
	"fmt"
	"strings"
)

// Quality is either Attack, Defend or Useful
type Quality struct {
	Type       string
	Name       string
	Dice       *DiePool
	Level      int
	Capacities []*Capacity
	Modifiers  []*Modifier
	CostPerDie int
}

func (q Quality) String() string {
	text := fmt.Sprintf("%s ", q.Type)

	// Add formatting for additional levels of Quality
	if q.Level > 0 {
		text += fmt.Sprintf(" +%d ", q.Level)
	}

	text += fmt.Sprintf("(%s) (%d/die): ",
		q.Name,
		q.CostPerDie)

	if len(q.Capacities) > 0 {

		text += fmt.Sprint("Capacities:")

		for _, c := range q.Capacities {
			text += fmt.Sprintf(" %s", c)
		}
	}

	if len(q.Modifiers) > 0 {
		text += fmt.Sprint("; Extras & Flaws:")

		for _, m := range q.Modifiers {
			if m.CostPerLevel > 0 {
				text += fmt.Sprintf(" %s,", m)
			}
		}

		for _, m := range q.Modifiers {
			if m.CostPerLevel < 0 {
				text += fmt.Sprintf(" %s,", m)
			}
		}
	}

	text = strings.TrimSuffix(text, ",")
	return text
}

// FormatDiePool returns a die string
func (q *Quality) FormatDiePool(actions int) string {

	for _, m := range q.Modifiers {
		if m.Name == "Spray" {
			q.Dice.Spray = m.Level
		}

		if m.Name == "Go First" {
			q.Dice.GoFirst = m.Level
		}
	}

	text := fmt.Sprintf("%dac+%dd+%dhd+%dwd+%dgf+%dsp",
		actions,
		q.Dice.Normal,
		q.Dice.Hard,
		q.Dice.Wiggle,
		q.Dice.GoFirst,
		q.Dice.Spray)

	return text
}

// NewQuality generates a new empty Quality
func NewQuality(t string) *Quality {

	q := new(Quality)

	q.Type = t
	q.Name = ""
	q.CostPerDie = 2
	q.Level = 0
	q.Capacities = []*Capacity{}
	q.Modifiers = []*Modifier{}

	// Take user input

	return q
}

// CalculateCost determines the cost of a Power Quality
// Called from Power.PowerCost()
func (q *Quality) CalculateCost(b int) {

	b += q.Level

	for _, m := range q.Modifiers {
		if m.RequiresLevel {
			b += m.CostPerLevel * m.Level
		} else {
			b += m.CostPerLevel
		}
	}
	if b < 0 {
		// No negative costs allowed
		b = 0
	}
	q.CostPerDie = b
}
