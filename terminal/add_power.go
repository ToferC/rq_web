package terminal

import (
	"fmt"
	"strconv"

	"github.com/go-pg/pg"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

func AddPower(db *pg.DB, c *oneroll.Character) {

	if len(c.Powers) == 0 {
		c.Powers = map[string]*oneroll.Power{}
	}

AddPowerLoop:
	for true {

		p := oneroll.Power{}

		fmt.Println("\nAdding a Power")

		if len(c.Archetype.Sources) < 1 {
			fmt.Println("\nYou need to have identified your Archetype, Sources and Permissions to purchase powers.")
			fmt.Println("\nCreating your Archetype")
			AddArchtype(db, c)
		}

		answer := UserQuery("\nEnter the name of your Power or hit enter to exit: ")

		if answer == "" {
			break AddPowerLoop
		}

		p.Name = answer

		err := ChooseQualities(&p)
		if err != nil {
			panic(err)
		}

		// Select Power Dice
		err = ChoosePowerDice(&p)
		if err != nil {
			panic(err)
		}

		// Calculate power capacities
		p.DeterminePowerCapacities()

		oneroll.UpdateCost(&p)

		fmt.Printf("Your Power:\n%s", p)

		for _, q := range p.Qualities {
			fmt.Println(q)
		}

		// Get user input for power Effect
		answer = UserQuery("\nDescribe your power's effect: ")

		p.Effect = answer

		c.Powers[p.Name] = &p

		fmt.Println(p)

		// Save character
		err = database.UpdateCharacter(db, c)
		if err != nil {
			panic(err)
		}
	} // End of AddPowerLoop
}

func ChooseQualities(p *oneroll.Power) error {

	var err error

ChooseQualitiesLoop:
	for true {

		fmt.Println("\nQualities:")

		qualities := []string{"Attack", "Defend", "Useful"}

		for _, q := range qualities {
			fmt.Printf("-- %s\n", q)
		}

		fmt.Printf("\nType the names of the Qualities you'd like to add one at a time. Hit Enter to move on.")

		answer := UserQuery("\nYour selection: ")

		if answer == "" {
			fmt.Println("Exiting.")
			break ChooseQualitiesLoop
		}

		validQuality := false

		for _, q := range qualities {
			if answer == q {
				validQuality = true
				break
			}
			validQuality = false
		}

		if !validQuality {
			fmt.Println("Not a valid Quality. Try again.")

		} else {

			q := new(oneroll.Quality)
			q.Type = answer

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

		} // End of Quality
	} // End of Quality Loop
	return err
}

// SelectQualityLevel gets user input for Quality Level and updates Quality
func SelectQualityLevel(q *oneroll.Quality) error {

	var err error

SelectQualityLevelLoop:
	for true {

		var num int
		var err error

		answer := UserQuery("Select additional levels for " + q.Type + " or hit enter to leave at base: ")
		if answer == "" {
			num = 0
		} else {
			num, err = strconv.Atoi(answer)
		}

		if err != nil || num < 0 {
			fmt.Println("Invalid value")
		} else {
			q.Level = num
			break SelectQualityLevelLoop
		}
	}
	return err
}

// SelectQualityLevel gets user input for Quality Level and updates Quality
func ChooseCapacities(q *oneroll.Quality) error {

	var err error
	capacities := map[string]int{
		"Mass":  25,
		"Range": 10,
		"Speed": 2,
		"Self":  0,
	}

	fmt.Println("You get one power capacity free per quality. Additional Capacities can be purchased through the Extra Capacity Extra.")

ChooseCapacitiesLoop:
	for true {

		c1 := new(oneroll.Capacity)

		fmt.Println("\nCapacities")

		for k, _ := range capacities {
			fmt.Printf("--%s\n", k)
		}

		fmt.Print("Current Capacities: ")
		for _, c := range q.Capacities {
			fmt.Printf("%s ", c.Type)
		}

		answer := UserQuery("\nSelect a Capacity for " + q.Type + " or hit enter to finish: ")

		if answer == "" {
			fmt.Println("Exiting.")
			break ChooseCapacitiesLoop
		}

		validCapacity := false

		for k := range capacities {
			if answer == k {
				validCapacity = true
				break
			} else {
				validCapacity = false
			}
		}

		if !validCapacity {
			fmt.Println("Not a valid Capacity. Try again.")

		} else {
			c1.Type = answer

			q.Capacities = append(q.Capacities, c1)
		}
	} // End ChooseCapacitiesLoop
	return err
}

func ChooseModifiers(q *oneroll.Quality) error {

	var err error

	// Add Power Capacity Mod if needed
	if len(q.Capacities) > 1 {
		tm := oneroll.Modifiers["Power Capacity"]
		tm.Level = len(q.Capacities) - 1
		q.Modifiers = append(q.Modifiers, &tm)
		fmt.Println("\n**Added Power Capacity Extra at level ", tm.Level)
	}

ExtrasOrFlawsLoop:
	for true {

		fmt.Println("\nExtras & Flaws:")

		answer := UserQuery(`
    1: Add Extra
    2: Add Flaw

    Or hit Enter to exit: `)

		if len(answer) == 0 {
			fmt.Println("Exiting")
			break ExtrasOrFlawsLoop
		}

		switch answer {
		case "1":
			fmt.Println("\nExtras:")

			for k, v := range oneroll.Modifiers {
				if v.CostPerLevel > 0 {
					fmt.Printf("-- %s (%dpts)\n", k, v.CostPerLevel)
				}
			}
		case "2":
			fmt.Println("\nFlaws:")

			for k, v := range oneroll.Modifiers {
				if v.CostPerLevel < 0 {
					fmt.Printf("-- %s (%dpts)\n", k, v.CostPerLevel)
				}
			}
		default:
			fmt.Println("Not a valid option. Please choose again")
		}

	ChooseModifiersLoop:
		for true {
			fmt.Printf("\nType the names of the Extras or Flaws you'd like to add one at a time. Hit Enter to finish adding.\n")

			fmt.Print("Current Extras & Flaws: ")
			for _, mod := range q.Modifiers {
				fmt.Printf("%s ", mod.Name)
			}

			answer = UserQuery("\nYour selection: ")

			if answer == "" {
				fmt.Println("Exiting.")
				break ChooseModifiersLoop
			}

			validModifier := false

			for _, k := range oneroll.Modifiers {
				if answer == k.Name {
					validModifier = true
					break
				}
				validModifier = false
			}

			if !validModifier {
				fmt.Println("Not a valid Modifier. Try again.")

			} else {

				// Create Modifiers
				m := oneroll.Modifiers[answer]

				// Add Info if required
				if m.RequiresInfo {
					answer := UserQuery("Describe the Extra or Flaw: ")

					m.Info = answer
				}

				// Add Level if required
				if m.RequiresLevel {
				LevelLoop:
					for true {
						answer := UserQuery("Enter the levels of your Extra of Flaw: ")
						num, err := strconv.Atoi(answer)

						if err != nil {
							fmt.Println("Invalid value")
						} else {
							m.Level = num
							break LevelLoop
						}
					} // End LevelLoop
				}
			} // End ChooseModifiersLoop

			tempMod := oneroll.Modifiers[answer]
			// Add the selected source to Quality.Modifiers
			q.Modifiers = append(q.Modifiers, &tempMod)
		}
	} // End ExtrasOrFlawsLoop
	return err
}

func ChoosePowerDice(p *oneroll.Power) error {

	var err error

	d := oneroll.DiePool{
		Normal: 0,
		Hard:   0,
		Wiggle: 0,
	}

	fmt.Println("Enter the values for your Power's dice pool.")

NormalDiceLoop:
	for true {
		answer := UserQuery("\nNormal Dice: ")
		num, err := strconv.Atoi(answer)

		// Response is a non-negative number
		if err != nil || num < 0 || num > 20 {
			fmt.Println("Invalid value")
		} else {
			d.Normal = num
			break NormalDiceLoop
		}
	} // End NormalDiceLoop

HardDiceLoop:
	for true {
		answer := UserQuery("\nHard Dice: ")
		num, err := strconv.Atoi(answer)

		// Response is a non-negative number
		if err != nil || num < 0 || num > 20 {
			fmt.Println("Invalid value")
		} else {
			d.Hard = num
			break HardDiceLoop
		}
	} // End HardDiceLoop

WiggleDiceLoop:
	for true {
		answer := UserQuery("\nWiggle Dice: ")
		num, err := strconv.Atoi(answer)

		// Response is a non-negative number
		if err != nil || num < 0 || num > 20 {
			fmt.Println("Invalid value")
		} else {
			d.Wiggle = num
			break WiggleDiceLoop
		}
	} // End WiggleDiceLoop

	p.Dice = &d

	return err
}
