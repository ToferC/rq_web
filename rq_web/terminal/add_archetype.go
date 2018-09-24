package terminal

import (
	"fmt"
	"os"

	"github.com/go-pg/pg"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/database"
)

func AddArchtype(db *pg.DB, c *oneroll.Character) {

	a := oneroll.Archetype{}

	fmt.Println("Adding an Archetype")

	answer := UserQuery("Enter the name of your Archetype or hit enter to exit: ")

	if answer == "" {
		os.Exit(3)
	}

	a.Type = answer

ChooseSourcesLoop:
	for true {

		fmt.Println("\nSources:")

		for k, v := range oneroll.Sources {
			fmt.Printf("-- %s (%dpts)\n", k, v.Cost)
		}

		fmt.Printf("\nType the names of the sources you'd like to add one at a time. Hit Enter to move on.")

		answer := UserQuery("\nYour selection: ")

		if answer == "" {
			fmt.Println("Exiting.")
			break ChooseSourcesLoop
		}

		validSource := false

		for k := range oneroll.Sources {
			if answer == k {
				validSource = true
				break
			}
			validSource = false
		}

		if !validSource {
			fmt.Println("Not a valid Source. Try again.")

		} else {
			// Add the selected source to Archetype.Sources
			tA := oneroll.Sources[answer]

			a.Sources = append(a.Sources, &tA)
		}
	} // End ChooseSourcesLoop

ChoosePermissionsLoop:
	for true {

		fmt.Println("\nPermissions:")

		for k, v := range oneroll.Permissions {
			fmt.Printf("-- %s (%dpts)\n", k, v.Cost)
		}

		fmt.Printf("\nType the names of the Permissions you'd like to add one at a time. Hit Enter to move on.")

		answer := UserQuery("\nYour selection: ")

		if answer == "" {
			fmt.Println("Exiting.")
			break ChoosePermissionsLoop
		}

		validPermission := false

		for k := range oneroll.Permissions {
			if answer == k {
				validPermission = true
				break
			}
			validPermission = false
		}

		if !validPermission {
			fmt.Println("Not a valid Permission. Try again.")

		} else {
			// Add the selected source to Archetype.Sources
			tP := oneroll.Permissions[answer]
			a.Permissions = append(a.Permissions, &tP)
		}
	} // End ChoosePermissionsLoop

ChooseIntrinsicsLoop:
	for true {

		fmt.Println("\nIntrinsics:")

		for k, v := range oneroll.Intrinsics {
			fmt.Printf("-- %s (%dpts)\n", k, v.Cost)
		}

		fmt.Printf("\nType the names of the Intrinsics you'd like to add one at a time. Hit Enter to move on.")

		answer := UserQuery("\nYour selection: ")

		if answer == "" {
			fmt.Println("Exiting.")
			break ChooseIntrinsicsLoop
		}

		validIntrinsic := false

		for k := range oneroll.Intrinsics {
			if answer == k {
				validIntrinsic = true
				break
			}
			validIntrinsic = false
		}

		if !validIntrinsic {
			fmt.Println("Not a valid Intrinsic. Try again.")

		} else {
			// Add the selected source to Archetype.Sources
			tI := oneroll.Intrinsics[answer]
			a.Intrinsics = append(a.Intrinsics, &tI)
		}
	} // End ChoosePermissionsLoop

	oneroll.UpdateCost(&a)
	c.Archetype = &a

	// Save character
	err := database.UpdateCharacter(db, c)
	if err != nil {
		panic(err)
	}
}
