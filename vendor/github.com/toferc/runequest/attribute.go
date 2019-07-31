package runequest

import (
	"fmt"
	"math"
)

// Attribute is a Character element that is based off other elements
type Attribute struct {
	Name            string
	MaxValue        int
	Value           int
	Base            int
	Updates         []*Update
	UserString      string
	Total           int
	Dice            int
	Max             int
	Min             int
	Text            string
	ExperienceCheck bool
}

// SetAttributes determines initial derived stats for the character
func (c *Character) SetAttributes() {

	var hp, hr, sz, db, dx, sd, mp *Attribute

	hp = c.DetermineHitPoints()
	hr = c.DetermineHealingRate()
	sz = c.DetermineSizStrikeRank()
	db = c.DetermineDamageBonus()
	sd = c.DetermineSpiritDamage()
	mp = c.DetermineMagicPoints()
	dx = c.DetermineDexStrikeRank()

	c.Attributes = map[string]*Attribute{
		"MP":    mp,
		"HP":    hp,
		"HR":    hr,
		"DB":    db,
		"SD":    sd,
		"DEXSR": dx,
		"SIZSR": sz,
	}
}

// UpdateAttributes totals base & value for Attribute
func (c *Character) UpdateAttributes() {
	for _, v := range c.Attributes {

		updates := 0

		for _, u := range v.Updates {
			updates += u.Value
		}

		v.Total = v.Base + v.Value + updates
	}
}

func (a *Attribute) String() string {

	var text string

	if a.Text == "" {
		text = fmt.Sprintf("%s: %d", a.Name, a.Total)
	} else {
		text = fmt.Sprintf("%s: %s", a.Name, a.Text)
	}
	return text
}

// DetermineMagicPoints calculates magic points for a Character
func (c *Character) DetermineMagicPoints() *Attribute {
	mp := &Attribute{
		Name: "Magic Points",
	}

	p := c.Statistics["POW"]
	p.UpdateStatistic()

	if p.Total > 0 {
		mp.Base = p.Total
		mp.Max = p.Total
		mp.MaxValue = p.Max
	} else {
		mp.Base = 0
		mp.Max = 0
	}

	return mp
}

// DetermineHitPoints calculates hit points for a Character
func (c *Character) DetermineHitPoints() *Attribute {

	hp := &Attribute{
		Name: "Hit Points",
	}

	var baseHP int

	s := c.Statistics["SIZ"]
	s.UpdateStatistic()

	siz := s.Total

	p := c.Statistics["POW"]
	p.UpdateStatistic()

	pow := p.Total

	con := c.Statistics["CON"]
	con.UpdateStatistic()

	if siz > 0 && con.Total > 0 {

		baseHP = con.Total

		switch {
		case siz < 5:
			baseHP -= 2
		case siz < 9:
			baseHP--
		case siz < 13:
		case siz < 17:
			baseHP++
		case siz < 21:
			baseHP += 2
		case siz < 25:
			baseHP += 3
		case siz > 24:
			f := float64(siz) - 24.0
			baseHP += int(math.Ceil(f/4)) + 4
		}

		switch {
		case pow == 0:
		case pow < 5:
			baseHP--
		case pow < 9:
		case pow < 13:
		case pow < 17:
		case pow < 21:
			baseHP++
		case pow < 25:
			baseHP += 2
		case pow > 24:
			f := float64(pow) - 24.0
			baseHP += int(math.Ceil(f/4)) + 3
		}
	}

	hp.Base = baseHP
	hp.Max = baseHP

	locHP := 0

	switch {
	case hp.Base == 0:
		locHP = 0
	case hp.Base < 7:
		locHP = 2
	case hp.Base < 10:
		locHP = 3
	case hp.Base < 13:
		locHP = 4
	case hp.Base < 16:
		locHP = 5
	case hp.Base < 19:
		locHP = 6
	case hp.Base < 22:
		locHP = 7
	case hp.Base > 21:
		f := float64(hp.Base) - 21.0
		locHP = int(math.Ceil(f/3)) + 7
	}

	for _, v := range c.HitLocations {
		v.Max = locHP + v.Base
		v.Min = -(2 * v.Max)
	}

	return hp
}

// DetermineHealingRate determines the Character's healingrate based on Con
func (c *Character) DetermineHealingRate() *Attribute {

	healingRate := &Attribute{
		Name: "Healing Rate",
	}

	con := c.Statistics["CON"]

	con.UpdateStatistic()

	tCon := con.Total

	switch {
	case tCon == 0:
		healingRate.Base = 0
	case tCon < 7:
		healingRate.Base = 1
	case tCon < 13:
		healingRate.Base = 2
	case tCon < 19:
		healingRate.Base = 3
	case tCon > 18:
		f := float64(healingRate.Base) - 18.0
		healingRate.Base = int(math.Ceil(f/6)) + 3
	}
	healingRate.Total = healingRate.Base + healingRate.Value
	return healingRate
}

// DetermineDamageBonus determines the Character's healingrate based on Con
func (c *Character) DetermineDamageBonus() *Attribute {

	damageBonus := &Attribute{
		Name: "Damage Bonus",
		Dice: 1,
	}

	str := c.Statistics["STR"]
	siz := c.Statistics["SIZ"]

	str.UpdateStatistic()
	siz.UpdateStatistic()

	db := siz.Total + str.Total

	switch {
	case db == 0:
		damageBonus.Base = 0
		damageBonus.Dice = 0
		damageBonus.Text = "N/A"
	case db < 13:
		damageBonus.Base = -4
		damageBonus.Text = "-1D4"
		damageBonus.Dice = 1
	case db < 25:
		damageBonus.Base = 0
		damageBonus.Dice = 0
		damageBonus.Text = ""
	case db < 33:
		damageBonus.Base = 4
		damageBonus.Dice = 1
		damageBonus.Text = "+1D4"
	case db < 41:
		damageBonus.Base = 6
		damageBonus.Dice = 1
		damageBonus.Text = "+1D6"
	case db < 57:
		damageBonus.Base = 6
		damageBonus.Dice = 2
		damageBonus.Text = "+2D6"
	case db > 56:
		damageBonus.Base = 6
		f := float64(db) - 56
		damageBonus.Dice = int(math.Ceil(f/16)) + 2
		damageBonus.Text = fmt.Sprintf("+%dD%d",
			damageBonus.Dice,
			damageBonus.Base,
		)
	}

	return damageBonus
}

// DetermineSpiritDamage determines the Character's healingrate based on Con
func (c *Character) DetermineSpiritDamage() *Attribute {

	damage := &Attribute{
		Name: "Spirit Damage",
		Dice: 1,
	}

	pow := c.Statistics["POW"]
	cha := c.Statistics["CHA"]

	pow.UpdateStatistic()
	cha.UpdateStatistic()

	db := pow.Total + cha.Total

	switch {
	case db == 0:
		damage.Base = 0
		damage.Text = "N/A"
	case db < 13:
		damage.Base = 3
		damage.Text = "1D3"
	case db < 25:
		damage.Base = 6
		damage.Text = "1D6"
	case db < 33:
		damage.Base = 6
		damage.Value = 1
		damage.Text = "1D6+1"
	case db < 41:
		damage.Base = 6
		damage.Value = 3
		damage.Text = "1D6+3"
	case db < 57:
		damage.Base = 6
		damage.Dice = 2
		damage.Value = 3
		damage.Text = "2D6+3"
	case db > 56:
		damage.Base = 6
		f := float64(db) - 56.0
		damage.Dice = int(math.Ceil(f/16)) + 2
		damage.Value = int(math.Ceil(f/16)) + 3
		damage.Text = fmt.Sprintf("%dD%d+%d",
			damage.Dice,
			damage.Base,
			damage.Value,
		)
	}
	return damage
}

// DetermineDexStrikeRank determines the Character's healingrate based on Con
func (c *Character) DetermineDexStrikeRank() *Attribute {

	dexSR := &Attribute{
		Name: "DEX Strike Rank",
		Max:  5,
	}

	dex := c.Statistics["DEX"]

	dex.UpdateStatistic()

	switch {
	case dex.Total == 0:
		dexSR.Base = 0
		dexSR.Text = "N/A"
	case dex.Total < 6:
		dexSR.Base = 5
	case dex.Total < 9:
		dexSR.Base = 4
	case dex.Total < 13:
		dexSR.Base = 3
	case dex.Total < 16:
		dexSR.Base = 2
	case dex.Total < 19:
		dexSR.Base = 1
	case dex.Total > 18:
		dexSR.Base = 0
	}
	return dexSR
}

// DetermineSizStrikeRank determines the Character's healingrate based on Con
func (c *Character) DetermineSizStrikeRank() *Attribute {

	sizSR := &Attribute{
		Name: "SIZ Strike Rank",
		Max:  5,
	}

	siz := c.Statistics["SIZ"]

	siz.UpdateStatistic()

	switch {
	case siz.Total == 0:
		sizSR.Base = 0
		sizSR.Text = "N/A"
	case siz.Total < 7:
		sizSR.Base = 3
	case siz.Total < 15:
		sizSR.Base = 2
	case siz.Total < 22:
		sizSR.Base = 1
	case siz.Total > 21:
		sizSR.Base = 0
	}
	return sizSR
}
