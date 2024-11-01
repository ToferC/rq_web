package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
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

	users, _ := database.ListUsers(db)

	for _, u := range users {

		if strings.Replace(u.UserName, " ", "", -1) == "" {
			fmt.Println("Found spam")
			database.DeleteUser(db, u.ID)
		}
	}
}
