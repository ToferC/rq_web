package terminal

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

func SelectPower(db *pg.DB, c *oneroll.Character) {

	fmt.Println("Choose a Power")

ModifyPowerLoop:
	for true {

		fmt.Println("Your Powers:")
		for _, p := range c.Powers {
			fmt.Printf("--%s (%dpts)\n", p.Name, p.Cost)
		}

		answer := UserQuery("\nType the name of the power to modify or hit Enter to exit: ")

		if answer == "" {
			fmt.Println("Exiting.")
			break ModifyPowerLoop
		}

		validPower := false

		for k := range c.Powers {
			if answer == k {
				validPower = true
				break
			}
			validPower = false
		}

		if !validPower {
			fmt.Println("Not a Power. Try again.")

		} else {

			p := c.Powers[answer]

			choosePowerAction(db, c, p)
		}
	}
}

func choosePowerAction(db *pg.DB, c *oneroll.Character, p *oneroll.Power) {
ChoosePowerOptionLoop:
	for true {
		fmt.Println("Choose an action:")

		answer := UserQuery(`
        1: Modify Power
  			2: Delete Power

  Or hit Enter to exit: `)

		if len(answer) == 0 {
			fmt.Println("Exiting")
			break ChoosePowerOptionLoop
		}

		switch answer {
		case "1":
			modifyPower(db, c, p)
		case "2":
			deletePower(db, c, p)
		default:
			fmt.Println("Not a valid option. Please choose again")
		}
	}
}

func deletePower(db *pg.DB, c *oneroll.Character, p *oneroll.Power) {

	response := UserQuery("Are you sure you want to delete " + p.Name + " ? (Y/N)")

	if response == "Y" || response == "y" {
		delete(c.Powers, p.Name)
		fmt.Println("Deleted.")
	} else {
		fmt.Println("Delete aborted.")
	}

	// Save character
	err := database.UpdateCharacter(db, c)
	if err != nil {
		panic(err)
	}

	fmt.Println("Deleted.")
}

func modifyPower(db *pg.DB, c *oneroll.Character, p *oneroll.Power) {

	for _, q := range p.Qualities {

		answer := UserQuery("Would you like to modify the " + q.Type + " Quality? (Y/N): ")

		if answer == "Y" || answer == "y" {

			answer := UserQuery("Briefly (1-4 words) describe your power quality: ")

			q.Name = answer

			// Choose level for qualities
			err := SelectQualityLevel(q)
			if err != nil {
				panic(err)
			}

			// Choose Capacities
			err = ChooseCapacities(q)
			if err != nil {
				panic(err)
			}

			// Choose Modifiers
			err = ChooseModifiers(q)
			if err != nil {
				panic(err)
			}

			// Add the completed quality to Power
			p.Qualities = append(p.Qualities, q)

		} else {
			fmt.Println("Skipping quality.")
		}

		err := ChooseQualities(p)
		if err != nil {
			panic(err)
		}

		// Select Power Dice
		//err = ChooseDice(p)
		if err != nil {
			panic(err)
		}

		// Calculate power capacities
		//hs.DeterminePowerCapacities()

		// Get user input for power Effect
		answer = UserQuery("\nDescribe your power's effect: ")

		p.Effect = answer

		oneroll.UpdateCost(p)

		fmt.Println(p)

		// Save character
		err = database.UpdateCharacter(db, c)
		if err != nil {
			panic(err)
		}
	} // End of ModifyPowerLoop
}
