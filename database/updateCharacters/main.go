package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-pg/pg/v10"
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

	cms, err := database.ListAllCharacterModels(db)

	for _, cm := range cms {
		if cm.Random == true && strings.TrimSpace(cm.Author.UserName) == "" {
			fmt.Printf("Found and deleted: %s\n", cm.Character.Name)
			database.DeleteCharacterModel(db, cm.ID)
		}
	}
}
