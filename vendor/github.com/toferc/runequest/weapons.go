package runequest

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Weapon represents a Runequest weapon
type Weapon struct {
	Name      string
	STR       int
	DEX       int
	Damage    string
	STRDamage bool
	HP        int
	CurrentHP int
	ENC       string
	Length    float64
	SR        int
	Type      string
	Range     int
	Special   string
}

// BaseWeapons is an array of runequest weapons
var BaseWeapons = loadWeapons()

func translateDieCode(s string) *DieCode {
	// translates a string like 1d6+1 into a DieCode

	var dice, mod, max int
	dieCode := &DieCode{}

	str := strings.Split(s, "+")
	if str[1] != "" {
		mod, _ = strconv.Atoi(str[1])
	}

	diceString := strings.Split(str[0], "D")
	dice, _ = strconv.Atoi(diceString[0])
	max, _ = strconv.Atoi(diceString[1])

	dieCode.NumDice = dice
	dieCode.DiceMax = max
	dieCode.Modifier = mod

	return dieCode
}

func loadWeapons() []*Weapon {

	weapons := []*Weapon{}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	url := "https://raw.githubusercontent.com/ToferC/runequest/master/weapons.csv"

	data, err := readCSVFromURL(url)

	for _, record := range data {

		sHP, _ := strconv.Atoi(record[5])

		if record[1] == "melee" {
			sSR, err := strconv.Atoi(record[2])
			if err != nil {
				sSR = 0
			}

			weapons = append(weapons, &Weapon{
				Name:      record[0],
				Type:      "Melee",
				SR:        sSR,
				ENC:       record[4],
				STRDamage: true,
				Damage:    record[3],
				HP:        sHP,
				CurrentHP: sHP,
			})
		} else {

			r, _ := strconv.Atoi(record[6])
			weapons = append(weapons, &Weapon{
				Name:      record[0],
				Type:      "Ranged",
				Range:     r,
				ENC:       record[4],
				Damage:    record[3],
				HP:        sHP,
				CurrentHP: sHP,
			})
		}
	}
	return weapons
}
