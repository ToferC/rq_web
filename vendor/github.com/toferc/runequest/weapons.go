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
	MainSkill string
	Damage    string
	STRDamage bool
	Thrown    bool
	HP        int
	CurrentHP int
	ENC       string
	Length    float64
	SR        int
	Type      string
	Range     int
	Special   string
	Custom    bool
}

// BaseWeapons is an array of runequest weapons
var BaseWeapons = loadWeapons()

// TranslateDieCode makes a string like 1d6+1 into a DieCode
func TranslateDieCode(s string) *DieCode {

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

	for i, record := range data {
		fmt.Println(i)

		sHP, err := strconv.Atoi(record[5])
		if err != nil {
			fmt.Println(err)
			sHP = 0
		}
		
		mainSkill := record[len(record)-1]

		if record[1] == "melee" {
			sSR, err := strconv.Atoi(record[2])
			if err != nil {
				sSR = 0
			}

			weapons = append(weapons, &Weapon{
				Name:      record[0],
				Type:      "Melee",
				MainSkill: mainSkill,
				SR:        sSR,
				ENC:       record[4],
				STRDamage: true,
				Damage:    record[3],
				HP:        sHP,
				CurrentHP: sHP,
			})
		} else {

			throw := false

			if strings.Contains(record[0], "Thrown") {
				throw = true
			}

			r, _ := strconv.Atoi(record[6])
			weapons = append(weapons, &Weapon{
				Name:      record[0],
				Type:      "Ranged",
				MainSkill: mainSkill,
				Range:     r,
				ENC:       record[4],
				Damage:    record[3],
				Thrown:    throw,
				HP:        sHP,
				CurrentHP: sHP,
			})

		}
	}
	return weapons
}
