package runequest

import "fmt"

// FamilyMember tracks family history through previous character history
type FamilyMember struct {
	Name        string
	Homeland    string
	Occupation  string
	Relation    string
	Born        int
	Died        int
	Alive       bool
	Description string
}

// Event represents a full event in previous character history
type Event struct {
	Year        int
	Name        string
	Start       bool
	End         bool
	Description string
	Participant string
	// Character, Parent, Grandparent
	HomelandModifiers map[string]int
	// Map Homeland name to modifier on d20 roll
	OccupationModifiers map[string]int
	// Map Occupation name to modifier on d20 roll
	Results           []EventResult
	FollowingEvent    string
	FollowingEventMod int
	Slug              string
}

// EventResult is a specific die range of random results from previous
// Character history
type EventResult struct {
	Range                []int
	Description          string
	Skills               []Skill
	Passions             []Ability
	Lunars               int
	Reputation           int
	Equipment            string
	Lethal               bool
	Boon                 bool
	RandomDeath          bool
	ImmediateFollowEvent string
	ImmediateFollowMod   int
	NextFollowEvent      string
	NextFollowMod        int
}

// Boon represents a random benefit from an EventResult
type Boon struct {
	Range       []int
	Description string
	Skills      []Skill
	Passsions   []Ability
	Lunars      int
	Reputation  int
	Equipment   string
}

// RandomCauseOfDeath is a random extinction event for family
type RandomCauseOfDeath struct {
	Range       []int
	Description string
	Skills      []Skill
	Passsions   []Ability
	Lunars      int
	Reputation  int
	Equipment   string
	Lethal      bool
}

// DetermineHistory generates a character's lifepath history
func DetermineHistory(c *Character, e Event, m int) (string, bool) {

	var homeland, occupation string

	fmt.Println(e.Participant)

	switch {
	case e.Participant == "Grandparent":
		if !c.Grandparent.Alive {
			return "Grandparent deceased", false
		}
		homeland = c.Grandparent.Homeland
		occupation = c.Grandparent.Occupation

	case e.Participant == "Parent":
		if !c.Parent.Alive {
			return "Parent deceased", false
		}
		homeland = c.Parent.Homeland
		occupation = c.Parent.Occupation

		fmt.Println("Character")
		homeland = c.Homeland.Name
		occupation = c.Occupation.Name
	}

	hlBonus, ok := e.HomelandModifiers[homeland]
	if !ok {
		return "didn't participate", false
	}
	ocBonus, ok := e.OccupationModifiers[occupation]
	if !ok {
		ocBonus = 0
	}

	roll := RollDice(20, 1, hlBonus+ocBonus+m, 1)

	if roll > 20 {
		roll = 20
	}

	if roll < 1 {
		roll = 1
	}

	for _, r := range e.Results {

		if IsInIntArray(r.Range, roll) {
			fmt.Println(r.Description, e.End)

			c.Lunars += r.Lunars
			c.Abilities["Reputation"].CreationBonusValue += r.Reputation

			text := r.Description

			for _, s := range r.Skills {
				// Update skills
				text += "\nSkills: "

				c.Skills[s.Name].Updates = append(c.Skills[s.Name].Updates, s.Updates[0])
				text += fmt.Sprintf("%s +%d%%", s.Name, s.Updates[0].Value)
			}

			for _, p := range r.Passions {
				// Update skills
				text += "\nPassions: "

				_, ok := c.Abilities[p.Name]

				if !ok {
					c.Abilities[p.Name] = &Ability{
						CoreString: p.CoreString,
						UserString: p.UserString,
						Base:       60,
						Updates:    []*Update{},
					}
				}

				fmt.Println(r.ImmediateFollowEvent, r.ImmediateFollowMod)

				c.Abilities[p.Name].Updates = append(c.Abilities[p.Name].Updates, p.Updates[0])
				text += fmt.Sprintf("%s +%d%%", p.Name, p.Updates[0].Value)
			}

			e.Description = text
			e.Results = []EventResult{r}
			c.History = append(c.History, &e)
			if e.End {
				return e.Name + " Ending", true
			}
		}
	}
	return e.Name + " Continuing", false
}
