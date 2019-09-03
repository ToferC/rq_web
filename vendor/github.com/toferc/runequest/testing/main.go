package main

import (
	"fmt"

	"github.com/toferc/runequest"
)

func main() {
	c := runequest.NewCharacter("Bob")

	c.Description = "Man"

	c.Grandparent.Homeland = "Esrolia"
	c.Grandparent.Occupation = "Hunter"

	mod := 0
	next := "1582_base"
	end := false

	for {

		event := runequest.PersonalHistoryEvents[next]

		fmt.Println(event)

		// Start history
		fmt.Println("Start History")
		result, _ := runequest.DetermineHistory(c, event, mod)

		fmt.Println(result)

		// Identify last EventResult

		last := c.History[len(c.History)-1]

		fmt.Println("LAST: ", last.Name)
		fmt.Println("NEXT: ", next)

		//immediate := last.Results[0].ImmediateFollowEvent
		// Check for immediate followup

		/*if immediate != "" {
			fmt.Println("New Immediate Event")

			for {
				// if immediate follow-up, go there
				e := runequest.PersonalHistoryEvents[immediate]
				mod := last.Results[0].ImmediateFollowMod
				r, end := runequest.DetermineHistory(c, e, mod)

				// Identify last EventResult
				last := c.History[len(c.History)-1]
				next := last.Results[0].NextFollowEvent

				// Check for immediate followup
				immediate := last.Results[0].ImmediateFollowEvent

				if immediate == "" {
					break
				}

				fmt.Println(next)
				fmt.Println(end, immediate, r)
			}
		}
		*/

		// if no immediate event or after immediate event
		fmt.Println("END: ", end)

		if last.End {
			break
		}
	}

	// Character done
	for k, v := range c.History {
		fmt.Println(k, v)
	}
}
