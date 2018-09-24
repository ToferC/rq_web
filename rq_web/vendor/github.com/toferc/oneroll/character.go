package oneroll

import (
	"fmt"
	"strings"
)

// Character represents a full character in the ORE game
type Character struct {
	ID           int64
	Name         string
	Setting      string
	Description  string
	Statistics   map[string]*Statistic
	StatMap      []string
	BaseWill     int
	Willpower    int
	Skills       map[string]*Skill
	Archetype    *Archetype
	HyperStats   map[string]*HyperStat
	HyperSkills  map[string]*HyperSkill
	Permissions  map[string]*Permission
	Powers       map[string]*Power
	Gear         string
	HitLocations map[string]*Location
	Passions     []*Passion
	Advantages   []*Advantage
	LocationMap  []string
	PointCost    int
	DetailedCost map[string]int
	InPlay       bool
	XP           int
	Updates      []*Update
}

// Update tracks live changes to Character
type Update struct {
	Date       string
	ChangeFrom string
	ChangeTo   string
	Cost       int
}

// Passion represents loyalties or drives
type Passion struct {
	Type        string
	Description string
	Value       int
}

// Display character
func (c *Character) String() string {

	text := fmt.Sprintf("\n%s (%d pts)\n", c.Name, c.PointCost)

	if c.Archetype.Type != "" {
		text += fmt.Sprint(c.Archetype)
	}

	text += "\n\nStats:\n"

	text += ShowSkills(c, false)

	text += fmt.Sprintf("\nBase Will: %d\n", c.BaseWill)
	text += fmt.Sprintf("Willpower: %d\n", c.Willpower)

	text += fmt.Sprintf("\nHit Locations:\n")

	for _, loc := range c.LocationMap {
		text += fmt.Sprintf("%s\n", c.HitLocations[loc])
	}

	if len(c.Archetype.Sources) > 0 {
		text += fmt.Sprintf("\nPowers:\n")

		for _, stat := range c.StatMap {
			s := c.Statistics[stat]
			if s.HyperStat != nil {
				text += fmt.Sprintf("\n%s\n", s.HyperStat)

				for _, q := range s.HyperStat.Qualities {
					text += fmt.Sprintf("%s\n", q)
				}

				if s.HyperStat.Effect != "" {
					text += fmt.Sprintf("Effect: %s", s.HyperStat.Effect)
				}

				if len(s.Modifiers) > 0 {
					cost := 0
					text += fmt.Sprintf("+ added modifiers to main stat: ")
					for _, m := range s.Modifiers {
						text += fmt.Sprintf("%s (%d/die) (%dpts), ", m.Name, m.Cost, m.Cost*SumDice(s.Dice))
						cost += m.Cost * SumDice(s.Dice)
					}
					text = strings.TrimSuffix(text, ",")
					text += fmt.Sprintf("(%dpts)", cost)
				}
				text += fmt.Sprint("\n\n")
			}
		}

		for _, s := range c.Skills {
			if s.HyperSkill != nil {
				text += fmt.Sprintf("\n%s\n", s.HyperSkill)

				for _, q := range s.HyperSkill.Qualities {
					text += fmt.Sprintf("%s\n", q)
				}

				if s.HyperSkill.Effect != "" {
					text += fmt.Sprintf("Effect: %s", s.HyperSkill.Effect)
				}

				if len(s.Modifiers) > 0 {
					cost := 0
					text += fmt.Sprintf("+ added modifiers to main skill: ")
					for _, m := range s.Modifiers {
						text += fmt.Sprintf("%s (%d/die) (%dpts), ", m.Name, m.Cost, m.Cost*SumDice(s.Dice))
						cost += m.Cost * SumDice(s.Dice)
					}
					text = strings.TrimSuffix(text, ",")
					text += fmt.Sprintf("(%dpts)", cost)
				}
				text += fmt.Sprint("\n")
			}
		}

		for _, p := range c.Powers {
			text += fmt.Sprintf("\n%s", p)

			for _, q := range p.Qualities {
				text += fmt.Sprintln(q)
			}

			if p.Effect != "" {
				text += fmt.Sprintf("Effect: %s\n", p.Effect)
			}
		}
	}
	return text
}

// CalculateCost updates the character and sums
// total costs of all character elements. Call this on each character update
func (c *Character) CalculateCost() {

	if !c.InPlay {

		var statsCost, skillsCost, powerCost int
		var archetypeCost, baseWillCost, willpowerCost int
		var advantageCost int

		if c.Setting != "RE" {
			if len(c.Archetype.Sources) > 0 {
				UpdateCost(c.Archetype)
				archetypeCost += c.Archetype.Cost
			}
		}

		for _, stat := range c.Statistics {
			UpdateCost(stat)
			statsCost += stat.Cost

			if stat.HyperStat != nil {
				UpdateCost(stat.HyperStat)
				powerCost += stat.HyperStat.Cost
			}
		}

		for _, skill := range c.Skills {
			UpdateCost(skill)
			skillsCost += skill.Cost

			if skill.HyperSkill != nil {
				UpdateCost(skill.HyperSkill)
				powerCost += skill.HyperSkill.Cost
			}
		}

		for _, power := range c.Powers {
			// Determine power capacities
			power.DeterminePowerCapacities()
			UpdateCost(power)
			powerCost += power.Cost
		}

		for _, advantage := range c.Advantages {
			if advantage.RequiresLevel {
				advantageCost += advantage.Cost * advantage.Level
			} else {
				advantageCost += advantage.Cost
			}
		}

		// Update BaseWill automaticallly if Character isn't in play

		calcBaseWill := 0

		if c.Setting != "RE" {
			for _, stat := range c.Statistics {
				if stat.EffectsWill {
					calcBaseWill += SumDice(stat.Dice)
					if stat.HyperStat != nil {
						calcBaseWill += SumDice(stat.HyperStat.Dice)
					}
				}
			}

			if c.BaseWill == 0 {
				// Auto-calculate base costs and levels for base character
				c.BaseWill = calcBaseWill
				c.Willpower = c.BaseWill
			}

			baseWillCost += 3 * (c.BaseWill - calcBaseWill)
			willpowerCost += c.Willpower - c.BaseWill

		}
		c.PointCost = archetypeCost + statsCost + skillsCost + powerCost + advantageCost + willpowerCost + baseWillCost

		c.DetailedCost = map[string]int{
			"archetype":  archetypeCost,
			"stats":      statsCost,
			"skills":     skillsCost,
			"powers":     powerCost,
			"advantages": advantageCost,
			"willpower":  willpowerCost,
			"basewill":   baseWillCost,
		}
	}
	//What happens when Character is in play
}
