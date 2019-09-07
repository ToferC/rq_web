package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-pg/pg"
	"github.com/toferc/rq_web/database"
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

	ocs, _ := database.ListOccupationModels(db)

	for _, oc := range ocs {

		equip := oc.Occupation.Equipment

		max := 6

		if len(equip) < 6 {
			max = len(equip)
		}
		for _, e := range equip[:max] {
			target := match(e)
			re := regexp.MustCompile("[0-9]+")
			armor := re.FindString(target)
			fmt.Println(armor)
			aVal, err := strconv.Atoi(armor)
			if err != nil {
				aVal = 0
			}
			oc.Occupation.GenericArmor = aVal
			break

		}
		database.UpdateOccupationModel(db, oc)
	}

}

func match(s string) string {
	i := strings.Index(s, "(")
	if i >= 0 {
		j := strings.Index(s[1:], ")")
		if j >= 0 {
			return s[i+1 : j-1]
		}
	}
	return ""
}
