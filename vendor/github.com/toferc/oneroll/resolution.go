package oneroll

import (
	"fmt"
	"sort"
)

// DieStringFormatter is an interface to generate a die string
// from a Stat, Power or Skill
type formatter interface {
	formatDiePool()
}

// OpposedRoll determines the results of an opposed roll between two or more actors
func OpposedRoll(rolls ...*Roll) []Match {

	fmt.Println("Opposed Roll Resolution")

	var results []Match

	for _, r := range rolls {

		fmt.Printf("Actor: %s, Action: %s, GoFirst: %d, Spray: %d, Wiggle Dice: %dwd\n",
			r.Actor.Name,
			r.Action,
			r.DiePool.GoFirst,
			r.DiePool.Spray,
			r.Wiggles,
		)

		for _, m := range r.Matches {
			results = append(results, m)
		}
		sort.Sort(ByWidthHeight(results))
	}
	return results
}

// PrintOpposed sorts actions by width and displays
func PrintOpposed(results []Match) {
	fmt.Println("***Resolution***")

	for i, m := range results {
		fmt.Printf("***ACTION %d: Actor: %s, Match: %dx%d, Initiative: %dx%d\n",
			i+1,
			m.Actor.Name,
			m.Height, m.Width,
			m.Height, m.Initiative,
		)
	}
}

// Resolve takes a slice of Rolls and determines outcomes
// This should probably be part of a Combat function
func Resolve(rolls ...*Roll) {

	fmt.Println("Opposed Roll Resolution")

	var results []Match

	// Initialize map to count actions
	var actors map[string]int
	var actions map[string][]string

	for _, r := range rolls {

		name := r.Actor.Name

		// Declarations

		// Sort by r.Actor.Sense

		// Track number of actions per actor
		actors[name] = r.NumActions

		// Declare actions for each action taken

	ChooseAction:
		for actors[r.Actor.Name] > 0 {
			fmt.Printf("Declare action %d of %d for %s.",
				actors[r.Actor.Name],
				r.NumActions,
				r.Actor.Name)

			answer := UserQuery(`
			1: Attack
			2: Defend
			3: Useful

			Choice: `)

			if answer == "" {
				fmt.Println("Exiting")
				break ChooseAction
			}

			switch answer {
			case "1":
				actions[name] = append(actions[name], "attack")
				actors[name]--
			case "2":
				actions[name] = append(actions[name], "defend")
				actors[name]--
			case "3":
				actions[name] = append(actions[name], "useful")
				actors[name]--
			default:
				fmt.Println("Not a valid option. Please choose again")
			}
		}

		fmt.Printf("Actor: %s, Action: %s, GoFirst: %d, Spray: %d, Wiggle Dice: %dwd\n",
			r.Actor.Name,
			r.Action,
			r.DiePool.GoFirst,
			r.DiePool.Spray,
			r.Wiggles,
		)

		for _, m := range r.Matches {
			results = append(results, m)
		}
		sort.Sort(ByWidthHeight(results))
	}

	//for _, m := range results {
	// In initiative order, let actors allocate their Matches
	// Attack = wound opponent and knock a die from their highest matche
	// Defend = gobble attacks against the actor
	// Useful = do something else
}
