package terminal

import (
	"fmt"
	"strconv"

	"github.com/fatih/structs"
	"github.com/go-pg/pg"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

// CreateCharacter takes terminal user input and saves to DB
func CreateCharacter(db *pg.DB) *oneroll.Character {

	name := UserQuery("What is the character's name? ")

	c := oneroll.NewSRCharacter(name)

	m := structs.Map(c)
	m["Name"] = name

	// Add statistics

	fmt.Println("\nAdding stats and skills.")

	fmt.Println("Enter normal die values (max 5) for:")

	for _, stat := range c.StatMap {
		s := c.Statistics[stat]

	StatsLoop:
		for true {
			answer := UserQuery("\n" + s.Name + ": ")
			num, err := strconv.Atoi(answer)

			if err != nil || num < 1 || num > 5 {
				fmt.Println("Invalid value")
			} else {
				s.Dice.Normal = num
				break StatsLoop
			}
		}

		for k, v := range c.Skills {
			if v.LinkStat.Name == s.Name {

			SkillsLoop:
				for true {

					str := fmt.Sprintf("-- %s: ", k)
					answer := UserQuery(str)
					num, err := strconv.Atoi(answer)

					if err != nil || num < 0 || num > 5 {
						fmt.Println("Invalid value")
					} else {
						c.Skills[k].Dice.Normal = num
						break SkillsLoop
					}
				}
			}
		}
	}

	oneroll.UpdateCost(c)

	fmt.Println(c)

	// Save character
	err := database.SaveCharacter(db, c)
	if err != nil {
		panic(err)
	}

	return c
}
