package main

import (
	"fmt"

	"github.com/toferc/runequest"
)

func main() {
	c := runequest.NewCharacter("Bob")

	c.Description = "Man"
	fmt.Println(c)

}
