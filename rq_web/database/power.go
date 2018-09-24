package database

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/toferc/oneroll"
)

// SavePower saves a Power to the DB
func SavePower(db *pg.DB, p *oneroll.Power) error {

	oneroll.UpdateCost(p)

	// Save character in Database
	_, err := db.Model(p).
		OnConflict("(id) DO UPDATE").
		Set("name = ?name").
		Insert(p)
	if err != nil {
		panic(err)
	}
	return err
}

func UpdatePower(db *pg.DB, p *oneroll.Power) error {

	oneroll.UpdateCost(p)

	err := db.Update(p)
	if err != nil {
		panic(err)
	}
	return err
}

// ListPowers queries Power names and add to slice
func ListPowers(db *pg.DB) (map[string]oneroll.Power, error) {
	var pows []oneroll.Power

	_, err := db.Query(&pows, `SELECT * FROM powers`)

	if err != nil {
		panic(err)
	}

	powMap := map[string]oneroll.Power{}

	// Create Map
	for i, p := range pows {
		powMap[oneroll.ToSnakeCase(p.Name)] = p
		fmt.Println(i, p.Name)
	}
	return powMap, nil
}

// LoadPower loads a single power from the DB by name
func LoadPower(db *pg.DB, name string) (*oneroll.Power, error) {
	// Select user by Primary Key
	pow := new(oneroll.Power)
	err := db.Model(pow).
		Where("Name = ?", name).
		Limit(1).
		Select()

	if err != nil {
		panic(err)
	}

	fmt.Println("Power loaded From DB")
	return pow, nil
}

// PKLoadPower loads a single power from the DB by pk
func PKLoadPower(db *pg.DB, pk int64) (*oneroll.Power, error) {
	// Select user by Primary Key
	pow := &oneroll.Power{ID: pk}
	err := db.Select(pow)

	if err != nil {
		return &oneroll.Power{Name: "New"}, err
	}

	fmt.Println("Power loaded From DB")
	return pow, nil
}

// DeletePower deletes a single power from DB by ID
func DeletePower(db *pg.DB, pk int64) error {

	pow := oneroll.Power{ID: pk}

	fmt.Println("Deleting power...")

	err := db.Delete(&pow)

	return err
}
