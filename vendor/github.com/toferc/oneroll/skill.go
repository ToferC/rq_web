package oneroll

import (
	"fmt"
	"sort"
)

// Skill represents specific training
type Skill struct {
	Name           string
	Quality        *Quality
	LinkStat       *Statistic
	Dice           *DiePool
	Narrow         bool
	Flexible       bool
	Influence      bool
	ReqSpec        bool
	Specialization string
	HyperSkill     *HyperSkill
	Modifiers      []*Modifier
	Free           bool
	CostPerDie     int
	Cost           int
}

// HyperSkill is a modified version of a regular Skill
type HyperSkill struct {
	Name       string
	Qualities  []*Quality
	Dice       *DiePool
	Effect     string
	Apply      bool
	CostPerDie int
	Cost       int
}

func (s Skill) String() string {

	td := ReturnDice(&s)

	text := ""

	if s.HyperSkill != nil {
		text += fmt.Sprintf("%s* ", s.Name)
	} else {
		text = fmt.Sprintf("%s ", s.Name)
	}

	if s.ReqSpec {
		text += fmt.Sprintf("[%s] ", s.Specialization)
	}

	if s.Narrow || s.Flexible || s.Influence {
		text += "("
		if s.Narrow {
			text += "N"
		}
		if s.Flexible {
			text += "F"
		}
		if s.Influence {
			text += "I"
		}
		text += ") "
	}

	text += fmt.Sprintf("%s", td)

	return text
}

func (hs HyperSkill) String() string {
	text := fmt.Sprintf("%s %s (", hs.Name, hs.Dice)

	for _, q := range hs.Qualities {
		text += fmt.Sprintf("%s", string(q.Type[0]))
		if q.Level > 0 {
			text += fmt.Sprintf("+%d", q.Level)
		}
	}

	text += fmt.Sprintf(") [%d/die] %dpts",
		hs.CostPerDie,
		hs.Cost)

	return text
}

// getDiePool returns a diepool based on a Skill and it's associated HyperSkill
func (s *Skill) getDiePool() *DiePool {

	td := &DiePool{}

	if s.HyperSkill != nil {

		td.Normal = s.Dice.Normal + s.HyperSkill.Dice.Normal
		td.Hard = s.Dice.Hard + s.HyperSkill.Dice.Hard
		td.Wiggle = s.Dice.Wiggle + s.HyperSkill.Dice.Wiggle
		td.Expert = s.HyperSkill.Dice.Expert

		for _, q := range s.HyperSkill.Qualities {
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
func (s *Skill) FormatDiePool(actions int) string {

	skill := ReturnDice(s)
	stat := ReturnDice(s.LinkStat)

	normal := stat.Normal + skill.Normal
	hard := stat.Hard + skill.Hard
	expert := skill.Expert
	wiggle := stat.Wiggle + skill.Wiggle
	goFirst := Max(stat.GoFirst, skill.GoFirst)
	spray := Max(stat.Spray, skill.Spray)

	text := fmt.Sprintf("%dac+%dd+%dhd+%dwd+%dgf+%dsp+%ded",
		actions,
		normal,
		hard,
		wiggle,
		goFirst,
		spray,
		expert)

	return text
}

// ShowSkills shows skills grouped under stats
// all bool determines if all skills are shown or just the ones with dice in them.
func ShowSkills(c *Character, allSkills bool) string {

	var text string

	for _, s := range c.StatMap {
		stat := c.Statistics[s]
		text += fmt.Sprintf("%s\n", stat)

		// Create Skill Mapping for stat
		skillMap := []string{}

		// Start with all skills
		for _, skill := range c.Skills {

			// Narrow down to only Skills with the right LinkStat
			if skill.LinkStat.Name == stat.Name {

				// Select all or only rated skills
				if allSkills {
					// We want all skills
					skillMap = append(skillMap, skill.Name)
				} else {
					// We only want rated skills
					if SkillRated(skill) {
						skillMap = append(skillMap, skill.Name)
					}
				}
			}
		}
		// Sort the map of Skills in Alphabetical order
		sort.Strings(skillMap)
		for _, skill := range skillMap {
			text += fmt.Sprintf("-- %s\n", c.Skills[skill])
		}

	}
	return text
}

// CalculateCost determines the cost of a Skill
// Called from Character.CalculateCharacterCost()
func (s *Skill) CalculateCost() {

	var b int

	switch {
	case s.Free:
		b = 0
	default:
		b = 2
	}

	if s.Narrow {
		b--
	}

	if s.Flexible {
		b++
	}

	if s.Influence {
		b++
	}

	s.CostPerDie = b

	// Modifier Cost
	mc := 0

	// Add cost for HyperSkill levels applied to Stat
	if s.HyperSkill != nil {
		for _, q := range s.HyperSkill.Qualities {
			if q.Level > 0 {
				b += q.Level
			}
		}
	}

	// Add costs for modifiers
	for _, m := range s.Modifiers {
		m.CalculateCost(0)
		if m.RequiresLevel {
			mc += m.CostPerLevel * m.Level
		} else {
			mc += m.CostPerLevel
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

	if s.Dice.Expert > 0 {
		total += 2
	}

	s.Cost = total
}

// CalculateCost generates and udpates the cost for HypeSKills
func (hs *HyperSkill) CalculateCost() {

	b := 1 // base of 1, but minimum of 1 Quality with minimum cost of 1

	for _, q := range hs.Qualities {

		// Add Power Capacity Modifier if needed
		if len(q.Capacities) > 1 {
			tm := Modifiers["Power Capacity"]
			tm.Level = len(q.Capacities) - 1
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
