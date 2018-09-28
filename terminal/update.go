package terminal

import (
	"fmt"
	"strconv"

	"github.com/go-pg/pg"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

// Update is a menu that provides access to Character upate options
func Update(db *pg.DB) {
	c, err := GetCharacter(db)

	if err != nil {
		panic(err)
	}

	oneroll.UpdateCost(c)

UpdateLoop:
	for true {
		fmt.Println("Choose an action:")

		answer := UserQuery(`
			1: Update Statistics
			2: Update Skills
			3: Add a Skill
			4: Delete a Skill
			5: Add an Archtype
			6: Add a Power
			7: Add Hyper-Stat
			9: Delete a Power
			8: Add Hyper-Skill
			10: Delete Hyper-Stats
			11: Delete Hyper-Skills

Or hit Enter to exit: `)

		if len(answer) == 0 {
			fmt.Println("Exiting")
			break UpdateLoop
		}

		switch answer {
		case "1":
			updateStat(db, c)
		case "2":
			updateSkills(db, c)
		case "3":
			AddSkill(db, c)
		case "4":
			deleteSkills(db, c)
		case "5":
			AddArchtype(db, c)
		case "6":
			AddPower(db, c)
		case "7":
			AddHyperStat(db, c)
		case "8":
			AddHyperSkill(db, c)
		case "9":
			deletePowers(db, c)
		case "10":
			deleteHyperStat(db, c)
		case "11":
			deleteHyperSkill(db, c)
		default:
			fmt.Println("Not a valid option. Please choose again")
		}
	}
}

func ChooseStatistic(c *oneroll.Character) *oneroll.Statistic {

	s := &oneroll.Statistic{}

UpdateStats:
	for true {

		fmt.Println("Your Stats")
		for _, s := range c.StatMap {
			fmt.Printf("-- %s\n", c.Statistics[s])
		}

		fmt.Printf("\nChoose the statistic to update or hit Enter to exit")

		answer := UserQuery("Your selection: ")

		if answer == "" {
			fmt.Println("Exiting.")
			break UpdateStats
		}

		validStat := false

		for k := range c.Statistics {
			if answer == k {
				validStat = true
				break
			}
			validStat = false
		}

		if !validStat {
			fmt.Println("Not a Stat. Try again.")

		} else {
			s = c.Statistics[answer]
			break UpdateStats
		}
	}
	return s
}

func updateStat(db *pg.DB, c *oneroll.Character) error {

	fmt.Println("Updating Statistics")

	s := ChooseStatistic(c)

	fmt.Println(s)

	fmt.Printf("%s has %d normal dice.\n", s.Name, s.Dice.Normal)
	nd := UserQuery("Please enter the new value: ")
	normal, err := strconv.Atoi(nd)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.Normal = normal
	}

	fmt.Printf("%s has %d hard dice.\n", s.Name, s.Dice.Hard)

	hd := UserQuery("Please enter the new value: ")
	hard, _ := strconv.Atoi(hd)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.Hard = hard
	}

	fmt.Printf("%s has %d wiggle dice.\n", s.Name, s.Dice.Wiggle)

	wd := UserQuery("Please enter the new value: ")
	wiggle, _ := strconv.Atoi(wd)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.Wiggle = wiggle
	}

	fmt.Println(c)

	// Update Linked Skills
	for _, skill := range c.Skills {
		if skill.LinkStat.Name == s.Name {
			skill.LinkStat = s
		}
	}

	// Save character
	err = database.UpdateCharacter(db, c)
	if err != nil {
		panic(err)
	}

	return err
}

func updateSkills(db *pg.DB, c *oneroll.Character) {

	fmt.Println("Updating Skills")

UpdateSkillsLoop:
	for true {

		fmt.Println(oneroll.ShowSkills(c, true))

		answer := UserQuery("\nType the name of the skill to update or hit Enter to exit: ")

		if answer == "" {
			fmt.Println("Exiting.")
			break UpdateSkillsLoop
		}

		validSkill := false

		for k := range c.Skills {
			if answer == k {
				validSkill = true
				break
			}
			validSkill = false
		}

		if !validSkill {
			fmt.Println("Not a skill. Try again.")

		} else {

			targetSkill := c.Skills[answer]

			err := updateSkill(db, targetSkill, c)
			if err != nil {
				panic(err)
			}
			fmt.Println("Updated. Choose another skill or hit Enter to exit.")
		}
	}
}

func updateSkill(db *pg.DB, s *oneroll.Skill, c *oneroll.Character) error {

	fmt.Println(s)

	if s.ReqSpec {
		fmt.Println("Current specialization is ", s.Specialization)
		spec := UserQuery("Please enter a new specialization or hit Enter to keep the current one: ")
		if len(spec) > 0 {
			s.Specialization = spec
		}
	}

	fmt.Printf("%s has %d normal dice.\n", s.Name, s.Dice.Normal)
	nd := UserQuery("Please enter the new value: ")
	normal, err := strconv.Atoi(nd)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.Normal = normal
	}

	fmt.Printf("%s has %d hard dice.\n", s.Name, s.Dice.Hard)

	hd := UserQuery("Please enter the new value: ")
	hard, _ := strconv.Atoi(hd)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.Hard = hard
	}

	fmt.Printf("%s has %d wiggle dice.\n", s.Name, s.Dice.Wiggle)

	wd := UserQuery("Please enter the new value: ")
	wiggle, _ := strconv.Atoi(wd)

	if err != nil {
		fmt.Println("Invalid value")
	} else {
		s.Dice.Wiggle = wiggle
	}

	fmt.Println(c)

	// Save character
	err = database.UpdateCharacter(db, c)
	if err != nil {
		panic(err)
	}

	return err
}

func AddSkill(db *pg.DB, c *oneroll.Character) {

AddSkillLoop:
	for true {

		fmt.Println(oneroll.ShowSkills(c, true))

		fmt.Println("Adding a new skill")

		s := oneroll.Skill{
			Name: "",
			Dice: &oneroll.DiePool{
				Normal: 0,
				Hard:   0,
				Wiggle: 0,
			},
			ReqSpec:        false,
			Specialization: "",
		}

		// Get user input for new skill

		answer := UserQuery("Enter the name of the new skill or hit Enter to exit: ")

		if answer == "" {
			break AddSkillLoop
		}

		s.Name = answer

		stat := ChooseStatistic(c)
		s.LinkStat = stat
		fmt.Println("Updated.")

		sp := UserQuery("Does the skill have a specialization? (Y/N):")

		if sp == "Y" || sp == "y" {
			s.ReqSpec = true
		}

		if s.ReqSpec {
			spec := UserQuery("Enter your specialization: ")
			s.Specialization = spec
		}

		c.Skills[answer] = &s

		updateSkill(db, &s, c)
	}
}

func deleteSkills(db *pg.DB, c *oneroll.Character) {

	fmt.Println("Deleting Skills")

DeleteSkillLoop:
	for true {

		fmt.Println(oneroll.ShowSkills(c, false))

		answer := UserQuery("\nType the name of the skill to delete or hit Enter to exit: ")

		if answer == "" {
			fmt.Println("Exiting.")
			break DeleteSkillLoop
		}

		validSkill := false

		for k := range c.Skills {
			if answer == k {
				validSkill = true
				break
			}
			validSkill = false
		}

		if !validSkill {
			fmt.Println("Not a skill. Try again.")

		} else {

			response := UserQuery("Are you sure you want to delete " + answer + " ? (Y/N)")

			if response == "Y" || response == "y" {
				delete(c.Skills, answer)
				fmt.Println("Deleted.")
			} else {
				fmt.Println("Delete aborted.")
			}

			// Save character
			err := database.UpdateCharacter(db, c)
			if err != nil {
				panic(err)
			}

		}
		fmt.Println("Deleted.")
	}
}

func deletePowers(db *pg.DB, c *oneroll.Character) {

	fmt.Println("Deleting Powers")

DeletePowerLoop:
	for true {

		for _, p := range c.Powers {
			fmt.Printf("--%s (%d)\n",
				p.Name, p.Cost)
		}

		answer := UserQuery("\nType the name of the power to delete or hit Enter to exit: ")

		if answer == "" {
			fmt.Println("Exiting.")
			break DeletePowerLoop
		}

		validPower := false

		for _, k := range c.Powers {
			if answer == k.Name {
				validPower = true
				break
			}
			validPower = false
		}

		if !validPower {
			fmt.Println("Not a power. Try again.")

		} else {

			response := UserQuery("Are you sure you want to delete " + answer + " ? (Y/N)")

			if response == "Y" || response == "y" {
				delete(c.Powers, answer)
				fmt.Println("Deleted.")
			} else {
				fmt.Println("Delete aborted.")
			}

			// Save character
			err := database.UpdateCharacter(db, c)
			if err != nil {
				panic(err)
			}

		}
		fmt.Println("Deleted.")
	}
}

func deleteHyperStat(db *pg.DB, c *oneroll.Character) {

	fmt.Println("Deleting Hyper Stat")

DeletePowerLoop:
	for true {
		for _, s := range c.Statistics {
			if s.HyperStat != nil {
				fmt.Printf("--%s (%d)\n",
					s.HyperStat.Name, s.HyperStat.Cost)
			}
		}

		answer := UserQuery("\nType the name of the base Stat to delete its Hyper-Stat or hit Enter to exit: ")

		if answer == "" {
			fmt.Println("Exiting.")
			break DeletePowerLoop
		}

		validPower := false

		for _, k := range c.Statistics {
			if k.HyperStat != nil {
				if answer == k.Name {
					validPower = true
					break
				}
			}
			validPower = false
		}

		if !validPower {
			fmt.Println(answer + " does not have a Hyper-Stat. Try again.")

		} else {

			response := UserQuery("Are you sure you want to delete " + answer + "Hyper-Stat? (Y/N): ")

			if response == "Y" || response == "y" {
				c.Statistics[answer].HyperStat = nil
				fmt.Println("Deleted.")
			} else {
				fmt.Println("Delete aborted.")
			}

			// Save character
			err := database.UpdateCharacter(db, c)
			if err != nil {
				panic(err)
			}

		}
		fmt.Println("Deleted.")
	}
}

func deleteHyperSkill(db *pg.DB, c *oneroll.Character) {

	fmt.Println("Deleting Hyper-Skill")

DeletePowerLoop:
	for true {
		for _, s := range c.Skills {
			if s.HyperSkill != nil {
				fmt.Printf("--%s (%d)\n",
					s.HyperSkill.Name, s.HyperSkill.Cost)
			}
		}

		answer := UserQuery("\nType the name of the base Stat to delete its Hyper-Stat or hit Enter to exit: ")

		if answer == "" {
			fmt.Println("Exiting.")
			break DeletePowerLoop
		}

		validPower := false

		for _, k := range c.Skills {
			if k.HyperSkill != nil {
				if answer == k.Name {
					validPower = true
					break
				}
			}
			validPower = false
		}

		if !validPower {
			fmt.Println(answer + " does not have a Hyper-Skill. Try again.")

		} else {

			response := UserQuery("Are you sure you want to delete " + answer + "Hyper-Skill? (Y/N): ")

			if response == "Y" || response == "y" {
				c.Skills[answer].HyperSkill = nil
				fmt.Println("Deleted.")
			} else {
				fmt.Println("Delete aborted.")
			}

			// Save character
			err := database.UpdateCharacter(db, c)
			if err != nil {
				panic(err)
			}

		}
		fmt.Println("Deleted.")
	}
}
