package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-pg/pg/v10"
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

	fmt.Println("Enter username for password change")

	username := runequest.UserQuery("Enter user name (or hit Enter to quit): ")

	if username == "" {
		os.Exit(3)
	}

	user, err := database.LoadUser(db, username)
	if err != nil {
		fmt.Println("No such user")
		os.Exit(3)
	}

	password := runequest.UserQuery("Enter new password: ")
	password2 := runequest.UserQuery("Re-enter new password: ")

	hashedPassword, err := database.HashPassword(password)
	if err != nil {
		fmt.Println(err)
	}

	if password != password2 {
		fmt.Println("Passwords do not match")
		os.Exit(3)
	}

	user.Password = hashedPassword

	database.UpdateUser(db, user)

	fmt.Println(user)
	fmt.Printf("Username %s password changed to %s", user.UserName, password)

}
