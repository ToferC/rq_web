package main

import (
	"fmt"

	"github.com/toferc/runequest"
)

func main() {

	c := runequest.NewCharacter("Default")
	c.UpdateCharacter()
	for k, v := range c.Skills {
		v.UpdateSkill()
		fmt.Println(k, v)
	}
	hl := runequest.SortLocations(c.HitLocations)
	fmt.Println(hl)

}
