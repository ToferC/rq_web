package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-pg/pg"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

// Query and return a Character from DB
func Query(db *pg.DB) {

	c, err := GetCharacter(db)

	if err != nil {
		panic(err)
	}

	// Ensure costs and validators are up to date
	oneroll.UpdateCost(c)

QueryActionLoop:
	for true {
		fmt.Println("Would you like to update Statistics or Skills?")

		answer := UserQuery(`
		1: Make a Skill Roll
		2: Mark Damage
		3: Coming Soon...
		4: Coming Soon...

		Or hit Enter to exit: `)

		if len(answer) == 0 {
			fmt.Println("Exiting")
			break QueryActionLoop
		}

		switch answer {
		case "1":
			rollSkill(c)
		case "2":
			//markDamage(db, c)
		default:
			fmt.Println("Not a valid option. Please choose again")
		}
	}
}

func rollSkill(c *oneroll.Character) {

ChooseSkillLoop:
	for true {
		fmt.Println("\nCharacter Skills:")

		fmt.Println(oneroll.ShowSkills(c, true))

		skillroll := UserQuery("\nChoose a skill to roll or hit Enter to quit: ")

		if skillroll == "" {
			fmt.Println("Exiting.")
			break ChooseSkillLoop
		}

		validSkill := true

		for k := range c.Skills {
			if skillroll == k {
				validSkill = true
				break
			}
			validSkill = false
		}

		if !validSkill {
			fmt.Println("Not a skill. Try again.")
		} else {

			s := c.Skills[skillroll]

			ds := s.FormatDiePool(1)

			fmt.Printf("Rolling %s (w/ %s) for %s\n",
				s.Name,
				s.LinkStat.Name,
				c.Name)

			r := oneroll.Roll{
				Actor:  c,
				Action: "Act " + s.Name,
			}

			r.Resolve(ds)

			fmt.Println("Rolling!")
			fmt.Println(r)
		}
	}
}

// GetCharacter lists all characters in DB and asks the user to select one
func GetCharacter(db *pg.DB) (*oneroll.Character, error) {

	var name string
	var err error

	// Select character loop
SelectCharacterLoop:
	for true {
		// List all charcters in DB
		list, err := database.ListCharacters(db)
		if err != nil {
			panic(err)
		}

		// Get user input on which character to load
		name = UserQuery("Enter your character's name to load or hit Enter to quit: ")

		if name == "" {
			fmt.Println("Exiting.")
			break SelectCharacterLoop
		}

		validCharacter := true

		for _, n := range list {
			if name == n.Name {
				validCharacter = true
				break
			}
			validCharacter = false
		}

		if validCharacter == false {
			fmt.Println("Not a valid character. Try again.")
		} else {
			break
		}
	}
	c, err := database.LoadCharacter(db, name)
	if err != nil {
		panic(err)
	}
	fmt.Println(c)
	return c, err
}

// UserQuery creates and question and returns the User's input as a string
func UserQuery(q string) string {
	question := bufio.NewReader(os.Stdin)
	fmt.Print(q)
	r, _ := question.ReadString('\n')

	input := strings.Trim(r, " \n")

	return input
}
