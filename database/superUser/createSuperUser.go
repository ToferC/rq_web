package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-pg/pg"
	"github.com/toferc/oneroll"
	"github.com/toferc/rq_web/database"
	"github.com/toferc/rq_web/models"
)

var db *pg.DB

func init() {
	os.Setenv("DBUser", "chris")
	os.Setenv("DBPass", "12345")
	os.Setenv("DBName", "ore_engine")
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

	for {
		var username, email, password, password2 string

		fmt.Println("Create SuperUser for ORE Engine")

		username = oneroll.UserQuery("Enter user name (or hit Enter to quit): ")

		if username == "" {
			break
		}

		email = oneroll.UserQuery("Enter user email: ")
		password = oneroll.UserQuery("Enter password: ")
		password2 = oneroll.UserQuery("Re-enter password: ")

		hashedPassword, err := database.HashPassword(password)
		if err != nil {
			fmt.Println(err)
		}

		if password != password2 {
			fmt.Println("Passwords do not match")
			break
		}
		user := models.User{
			UserName: username,
			Email:    email,
			Password: hashedPassword,
			IsAdmin:  true,
		}
		database.SaveUser(db, &user)
		fmt.Println(user)
		fmt.Printf("Superuser %s created", user.UserName)
		break
	}

}
