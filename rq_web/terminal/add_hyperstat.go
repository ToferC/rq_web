package terminal

import (
	"fmt"
	"strconv"

	"github.com/go-pg/pg"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

func AddHyperStat(db *pg.DB, c *oneroll.Character) {

	fmt.Println("Select Stat to Add Hyper-Stat to:")

	stat := ChooseStatistic(c)

	err := HyperStatInput(db, stat, c)
	if err != nil {
		panic(err)
	}
	fmt.Println("Hyper-Stat created. Choose another stat or hit Enter to exit.")
}

func HyperStatInput(db *pg.DB, s *oneroll.Statistic, c *oneroll.Character) error {

AddHyperStatLoop:
	for true {

		fmt.Println("\nAdding a Hyper-Stat")

		if len(c.Archetype.Sources) < 1 {
			fmt.Println("\nYou need to have identified your Archetype, Sources and Permissions to purchase Hyper-Stats.")
			fmt.Println("\nCreating your Archetype")
			AddArchtype(db, c)
		}

		s.HyperStat = &oneroll.HyperStat{}

		hs := s.HyperStat

		// Set all Qualities for Hyper-Stat
		hs.Name = "Hyper-" + s.Name

		hs.Qualities = []*oneroll.Quality{}

		qualities := []string{"Attack", "Defend", "Useful"}

		for _, qs := range qualities {
			q := &oneroll.Quality{
				Type:       qs,
				Level:      1,
				CostPerDie: 0,
			}

			answer := UserQuery("Would you like to modify the " + qs + " Quality? (Y/N): ")

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
				hs.Qualities = append(hs.Qualities, q)

			} else {
				fmt.Println("Skipping quality.")
			}

		}

		err := ChooseAdditionalHyperStatQualities(hs)
		if err != nil {
			panic(err)
		}

		// Select Power Dice
		err = ChooseHyperStatDice(hs)
		if err != nil {
			panic(err)
		}

		fmt.Println(hs)

		// Get user input for power Effect
		answer := UserQuery("\nDescribe your power's effect: ")

		hs.Effect = answer

		answer = UserQuery("Would you like to add your modifiers to the  " + s.Name + " Stat? (Y/N): ")

		if answer == "Y" || answer == "y" {

			for _, q := range hs.Qualities {
				for _, m := range q.Modifiers {

					s.Modifiers = append(s.Modifiers, m)
				}
			}
		}

		oneroll.UpdateCost(hs)

		fmt.Println(hs)

		// Save character
		err = database.UpdateCharacter(db, c)
		if err != nil {
			panic(err)
		}
		break AddHyperStatLoop
	} // End of AddHyperStatLoop
	return nil
}

func ChooseAdditionalHyperStatQualities(p *oneroll.HyperStat) error {

	var err error

ChooseQualitiesLoop:
	for true {

		fmt.Println("Hyper-Stats start with all 3 Qualities")
		fmt.Println("You have the following Qualities:")
		for _, q := range p.Qualities {
			fmt.Println(q)
		}

		fmt.Println("\nQualities:")

		qualities := []string{"Attack", "Defend", "Useful"}

		for _, q := range qualities {
			fmt.Printf("-- %s\n", q)
		}

		fmt.Printf("\nType the names of any additional Qualities you'd like to add one at a time. Hit Enter to move on.")

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

func ChooseHyperStatDice(p *oneroll.HyperStat) error {

	var err error

	d := oneroll.DiePool{
		Normal: 0,
		Hard:   0,
		Wiggle: 0,
	}

	fmt.Println("Enter the values for your HyperStat's dice pool.")

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
