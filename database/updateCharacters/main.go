package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-pg/pg"
	"github.com/toferc/rq_web/database"
	"github.com/toferc/runequest"
)

var db *pg.DB

func init() {
	os.Setenv("DBUser", "chris")
	os.Setenv("DBPass", "12345")
	os.Setenv("DBName", "runequest")
}

func main() {

	if os.Getenv("ENVIRONMENT") == "production" {
		url, ok := os.LookupEnv("DATABASE_URL")

		if !ok {
			log.Fatalln("$DATABASE_URL is required")
		}

		options, err := pg.ParseURL(url)

		if err != nil {
			log.Fatalf("Connection error: %s", err.Error())
		}

		db = pg.Connect(options)
	} else {
		db = pg.Connect(&pg.Options{
			User:     os.Getenv("DBUser"),
			Password: os.Getenv("DBPass"),
			Database: os.Getenv("DBName"),
		})
	}

	defer db.Close()

	err := database.InitDB(db)
	if err != nil {
		panic(err)
	}

	fmt.Println(db)

	// AddSlug to cms
	cults, _ := database.ListCultModels(db)

	for _, cl := range cults {

		tempRuneSpells := []*runequest.Spell{}

		for _, v := range cl.Cult.RuneSpells {

			index, err := indexSpell(v.CoreString, runequest.RuneSpells)
			if err != nil {
				fmt.Println(err)
				tempRuneSpells = append(tempRuneSpells, v)
				continue
			}

			baseSpell := runequest.RuneSpells[index]

			s := &runequest.Spell{
				Name:       baseSpell.Name,
				CoreString: baseSpell.CoreString,
				UserString: v.UserString,
				Cost:       baseSpell.Cost,
				Runes:      baseSpell.Runes,
				Domain:     baseSpell.Domain,
			}

			s.GenerateName()
			tempRuneSpells = append(tempRuneSpells, s)
		}
		cl.Cult.RuneSpells = tempRuneSpells

		database.UpdateCultModel(db, cl)
	}
}

func indexSpell(str string, spells []runequest.Spell) (int, error) {

	err := errors.New("Spell Not Found")

	for i, spell := range spells {
		if str == spell.CoreString {
			return i, nil
		}
	}

	return 0, err
}
