package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-pg/pg"
	"github.com/toferc/rq_web/models"
)

// SaveUser saves a User to the DB
func SaveUser(db *pg.DB, u *models.User) error {

	// Save User in Database
	_, err := db.Model(u).
		OnConflict("(id) DO UPDATE").
		Set("user_name = ?user_name").
		Insert(u)
	if err != nil {
		panic(err)
	}
	return err
}

//UpdateUser updates user info
func UpdateUser(db *pg.DB, u *models.User) error {

	err := db.Update(u)
	if err != nil {
		panic(err)
	}
	return err
}

//LoadUser will check if the user exists in db and if the
//username password combination is valid
func LoadUser(db *pg.DB, username string) (*models.User, error) {
	user := new(models.User)

	if strings.TrimSpace(username) != "" {

		fmt.Println("Loading User " + username)
		err := db.Model(user).
			Where("user_name = ?", username).
			Limit(1).
			Select()
		if err != nil {
			fmt.Println(err)
			return new(models.User), err
		}
		return user, nil
	} else {
		return nil, fmt.Errorf("Incorrect Username")
	}
}

//ValidUser will check if the user exists in db and if the
//username password combination is valid
func ValidUser(db *pg.DB, username, password string) bool {
	user := new(models.User)
	err := db.Model(user).
		Where("user_name = ?", username).
		Limit(1).
		Select()
	if err != nil {
		return false
	}

	if CheckPasswordHash(password, user.Password) {
		return true
	}
	return false
}

// ListUsers queries User names and add to slice
func ListUsers(db *pg.DB) ([]*models.User, error) {
	var users []*models.User

	_, err := db.Query(&users, `SELECT * FROM Users`)

	if err != nil {
		panic(err)
	}

	// Print names and PK
	for _, n := range users {
		n.Characters = countUserCharacterModels(db, n.UserName)
	}
	return users, nil
}

// PaginateUsers queries open Character names and add to slice
func PaginateUsers(db *pg.DB, limit, offset int) ([]*models.User, error) {

	var users []*models.User

	err := db.Model(&users).
		Limit(limit).
		Offset(offset * limit).
		Order("created_at DESC").
		Select()

	if err != nil {
		log.Println(err)
	}

	for _, n := range users {
		n.Characters = countUserCharacterModels(db, n.UserName)
	}

	return users, nil
}

// PKLoadUser loads a single User from the DB by pk
func PKLoadUser(db *pg.DB, pk int64) (*models.User, error) {
	// Select user by Primary Key
	user := &models.User{ID: pk}
	err := db.Select(user)

	if err != nil {
		return &models.User{UserName: "New"}, err
	}

	fmt.Println("User loaded From DB")
	return user, nil
}

// DeleteUser deletes a single User from DB by ID
func DeleteUser(db *pg.DB, pk int64) error {

	user := models.User{ID: pk}

	fmt.Println("Deleting User...")

	err := db.Delete(&user)

	return err
}

// CreateGoogleUser adds a user to the database based on a google user context
func CreateGoogleUser(db *pg.DB, username, email string) (*models.User, error) {

	fmt.Println(username, "no password", email)

	u := models.User{
		UserName: username,
		Email:    email,
	}

	err := SaveUser(db, &u)
	if err != nil {
		return &models.User{UserName: "New"}, err
	}

	return &u, nil
}
