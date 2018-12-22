package runequest

import (
	"fmt"
	"sort"
)

// Creature represents a generic RPG Creature
type Creature struct {
	Name string
	Role string
	// Type of Creature
	Description string

	Age          int
	Affiliations string
	Abilities    map[string]*Ability
	// Passions and Reputation
	ElementalRunes map[string]*Ability
	// Elemental Runes
	PowerRunes      map[string]*Ability
	StatisticFrames map[string]*StatisticFrame
	Statistics      map[string]*Statistic
	Attributes      map[string]*Attribute
	CurrentHP       int
	CurrentMP       int
	CurrentRP       int
	Cults           map[string]int
	Move            int
	DerivedMap      []string
	Skills          map[string]*Skill
	SkillMap        []string
	SkillCategories map[string]*SkillCategory
	RuneSpells      map[string]*Spell
	SpiritMagic     map[string]*Spell
	Powers          map[string]*Power
	HitLocations    map[string]*HitLocation
	HitLocationMap  []string
	MeleeAttacks    map[string]*Attack
	RangedAttacks   map[string]*Attack
	Equipment       []string
}

// CreatureRoles is an array of options for Creature.Role
var CreatureRoles = []string{
	"Non-Player Creature",
	"Animal",
	"Chaos",
	"Creature",
	"Demon",
	"Elemental",
	"Spirit",
}

// TotalStatistics updates values for stats after being modified
func (c *Creature) TotalStatistics() {

	for _, s := range c.Statistics {

		s.UpdateStatistic()
	}
}

// UpdateCreature updates stats, runes and skills based on them
func (c *Creature) UpdateCreature() {
	c.TotalStatistics()
	c.DetermineSkillCategoryValues()

	// This can be optimized
	c.SetAttributes()
	c.UpdateAttributes()

	if len(c.HitLocations) == 0 {
		c.HitLocationMap = HitLocationMap
		c.HitLocations = Locations
	}
}

// DetermineSkillCategoryValues figures out base values for skill categories based on stats
func (c *Creature) DetermineSkillCategoryValues() {

	for _, sc := range CategoryOrder {

		c.SkillCategories[sc] = &SkillCategory{}

		c.SkillCategories[sc].Name = sc
		c.SkillCategories[sc].Value = 0
	}

	for k, sc := range RQSkillCategories {
		// For each category

		for _, sm := range sc {
			// For each modifier in a category

			// Identify the stat
			stat := c.Statistics[sm.statistic]
			stat.UpdateStatistic()

			// Match to SkillCategory
			s := c.SkillCategories[k]

			// Map against specific values
			switch {
			case stat.Total <= 4:
				s.Value += sm.values[4]
			case stat.Total <= 8:
				s.Value += sm.values[8]
			case stat.Total <= 12:
				s.Value += sm.values[12]
			case stat.Total <= 16:
				s.Value += sm.values[16]
			case stat.Total <= 20:
				s.Value += sm.values[20]
			case stat.Total > 20:
				if sm.values[20] > 0 {
					s.Value += sm.values[20] + ((stat.Total-20)/4)*5
				} else {
					s.Value += sm.values[20] - ((stat.Total-20)/4)*5
				}
			}
		}
	}
}

// SetAttributes determines initial derived stats for the Creature
func (c *Creature) SetAttributes() {

	mp := c.DetermineMagicPoints()
	hp := c.DetermineHitPoints()
	hr := c.DetermineHealingRate()
	db := c.DetermineDamageBonus()
	sd := c.DetermineSpiritDamage()
	dx := c.DetermineDexStrikeRank()
	sz := c.DetermineSizStrikeRank()

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
func (c *Creature) UpdateAttributes() {
	for _, v := range c.Attributes {

		updates := 0

		for _, u := range v.Updates {
			updates += u.Value
		}

		v.Total = v.Base + v.Value + updates
	}
}

func (c *Creature) String() string {

	text := c.Name
	text += fmt.Sprintf("\nAffiliations: %s", c.Affiliations)

	text += "\n\nStats:\n"
	for _, stat := range StatMap {
		text += fmt.Sprintf("%s\n", c.Statistics[stat])
	}

	text += "\nDerived Stats:\n"
	for _, ds := range c.Attributes {
		text += fmt.Sprintf("%s\n", ds)
	}

	text += "\nAbilities:"

	for _, at := range AbilityTypes {

		text += fmt.Sprintf("\n\n**%s**", at)

		for _, ability := range c.Abilities {

			if ability.Type == at {
				text += fmt.Sprintf("\n%s", ability)
			}
		}
	}

	text += "\nElemental Runes:"

	for _, ability := range c.ElementalRunes {

		text += fmt.Sprintf("\n%s", ability)
	}

	text += "\nPower Runes:"

	for _, ability := range c.ElementalRunes {

		text += fmt.Sprintf("\n%s", ability)
	}

	text += "\n\nSkills:"

	keys := make([]string, 0, len(c.Skills))
	for k := range c.Skills {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, co := range CategoryOrder {

		sc := c.SkillCategories[co]

		text += fmt.Sprintf("%s", sc)
		for _, skill := range keys {

			if c.Skills[skill].Category == sc.Name {
				text += fmt.Sprintf("\n%s", c.Skills[skill])
			}
		}

	}
	return text
}

// DetermineMagicPoints calculates magic points for a Creature
func (c *Creature) DetermineMagicPoints() *Attribute {
	mp := &Attribute{
		Name:     "Magic Points",
		MaxValue: 21,
	}

	p := c.Statistics["POW"]
	p.UpdateStatistic()

	mp.Base = p.Total
	mp.Max = p.Total

	return mp
}

// DetermineHitPoints calculates hit points for a Creature
func (c *Creature) DetermineHitPoints() *Attribute {

	hp := &Attribute{
		Name:     "Hit Points",
		MaxValue: 21,
	}

	s := c.Statistics["SIZ"]
	s.UpdateStatistic()

	siz := s.Total

	p := c.Statistics["POW"]
	p.UpdateStatistic()

	pow := p.Total

	fmt.Println("SIZ ", siz)
	fmt.Println("POW ", pow)

	con := c.Statistics["CON"]
	con.UpdateStatistic()

	fmt.Println("CON ", con.Total)

	baseHP := con.Total

	switch {
	case siz < 5:
		baseHP -= 2
		fmt.Println("-2")
	case siz < 9:
		baseHP--
		fmt.Println("-1")
	case siz < 13:
		fmt.Println("No change")
	case siz < 17:
		baseHP++
		fmt.Println("+1")
	case siz < 21:
		baseHP += 2
		fmt.Println("+2")
	case siz < 25:
		baseHP += 3
		fmt.Println("+3")
	case siz > 24:
		baseHP += ((siz - 24) / 4) + 4
	}

	switch {
	case pow < 5:
		baseHP--
	case pow < 9:
		fmt.Println("No change")
	case pow < 13:
		fmt.Println("No change")
	case pow < 17:
		fmt.Println("No change")
	case pow < 21:
		baseHP++
	case pow < 25:
		baseHP += 2
	case pow > 24:
		baseHP += ((pow - 24) / 4) + 3
	}

	hp.Base = baseHP
	hp.Max = baseHP

	locHP := 0

	switch {
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
		locHP = ((hp.Base - 21) / 3) + 7
	}

	for _, v := range c.HitLocations {
		v.Max = locHP + v.Base
		v.Min = -(2 * v.Max)
	}

	return hp
}

// DetermineHealingRate determines the Creature's healingrate based on Con
func (c *Creature) DetermineHealingRate() *Attribute {

	healingRate := &Attribute{
		Name: "Healing Rate",
		Max:  21,
	}

	con := c.Statistics["CON"]

	con.UpdateStatistic()
	tCon := con.Total

	switch {
	case tCon < 7:
		healingRate.Base = 1
	case tCon < 13:
		healingRate.Base = 2
	case tCon < 19:
		healingRate.Base = 3
	case tCon > 18:
		healingRate.Base = ((tCon - 18) / 6) + 3
	}
	healingRate.Total = healingRate.Base + healingRate.Value
	return healingRate
}

// DetermineDamageBonus determines the Creature's healingrate based on Con
func (c *Creature) DetermineDamageBonus() *Attribute {

	damageBonus := &Attribute{
		Name: "Damage Bonus",
		Max:  21,
		Dice: 1,
	}

	str := c.Statistics["STR"]
	siz := c.Statistics["SIZ"]

	str.UpdateStatistic()
	siz.UpdateStatistic()

	db := siz.Total + str.Total

	switch {
	case db < 13:
		damageBonus.Base = -4
		damageBonus.Text = "-1D4"
	case db < 25:
		damageBonus.Base = 0
		damageBonus.Text = "-"
	case db < 33:
		damageBonus.Base = 4
		damageBonus.Text = "+1D4"
	case db < 41:
		damageBonus.Base = 6
		damageBonus.Text = "+1D6"
	case db < 57:
		damageBonus.Base = 6
		damageBonus.Dice = 2
		damageBonus.Text = "+2D6"
	case db > 56:
		damageBonus.Base = 6
		damageBonus.Dice = ((db - 56) / 16) + 2
		damageBonus.Text = fmt.Sprintf("+%dD%d",
			damageBonus.Dice,
			damageBonus.Base,
		)
	}

	return damageBonus
}

// DetermineSpiritDamage determines the Creature's healingrate based on Con
func (c *Creature) DetermineSpiritDamage() *Attribute {

	damage := &Attribute{
		Name: "Spirit Damage",
		Max:  21,
		Dice: 1,
	}

	pow := c.Statistics["POW"]
	cha := c.Statistics["CHA"]

	pow.UpdateStatistic()
	cha.UpdateStatistic()

	db := pow.Total + cha.Total

	switch {
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
		damage.Dice = ((db - 56) / 16) + 2
		damage.Value = ((db - 56) / 16) + 3
		damage.Text = fmt.Sprintf("%dD%d+%d",
			damage.Dice,
			damage.Base,
			damage.Value,
		)
	}
	return damage
}

// DetermineDexStrikeRank determines the Creature's healingrate based on Con
func (c *Creature) DetermineDexStrikeRank() *Attribute {

	dexSR := &Attribute{
		Name: "DEX Strike Rank",
		Max:  5,
	}

	dex := c.Statistics["DEX"]

	dex.UpdateStatistic()

	switch {
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

// DetermineSizStrikeRank determines the Creature's healingrate based on Con
func (c *Creature) DetermineSizStrikeRank() *Attribute {

	sizSR := &Attribute{
		Name: "SIZ Strike Rank",
		Max:  5,
	}

	siz := c.Statistics["SIZ"]

	siz.UpdateStatistic()

	switch {
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
