package database

import (
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/toferc/rq_web/models"
)

// SaveEncounter saves a Encounter to the DB
func SaveEncounter(db *pg.DB, enc *models.Encounter) error {

	// Save encounter in Database
	_, err := db.Model(enc).
		OnConflict("(id) DO UPDATE").
		Set("name = ?name").
		Insert(enc)
	if err != nil {
		panic(err)
	}
	return err
}

// UpdateEncounter updates a runequest encounter
func UpdateEncounter(db *pg.DB, enc *models.Encounter) error {

	_, err := db.Model(enc).Update()
	if err != nil {
		panic(err)
	}
	return err
}

// ListEncounters queries Encounter names and add to slice
func ListEncounters(db *pg.DB) ([]*models.Encounter, error) {
	var encs []*models.Encounter

	_, err := db.Query(&encs, `SELECT * FROM encounters`)

	if err != nil {
		panic(err)
	}

	return encs, nil
}

// ListUserEncounters queries Encounter names and add to slice
func ListUserEncounters(db *pg.DB, username string) ([]*models.Encounter, error) {
	var encs []*models.Encounter

	_, err := db.Query(&encs, `SELECT * FROM encounters WHERE author ->> 'UserName' = ?`, username)

	if err != nil {
		panic(err)
	}

	return encs, nil
}

// PKLoadEncounter loads a single encounter from the DB by pk
func PKLoadEncounter(db *pg.DB, pk int64) (*models.Encounter, error) {
	// Select user by Primary Key
	enc := &models.Encounter{ID: pk}
	err := db.Model(enc).Select()

	if err != nil {
		fmt.Println(err)
		return &models.Encounter{}, err
	}

	fmt.Println("Encounter loaded From DB")
	return enc, nil
}

// SlugLoadEncounter loads a single encounter from the DB by pk
func SlugLoadEncounter(db *pg.DB, slug string) (*models.Encounter, error) {
	// Select user by Primary Key
	enc := &models.Encounter{}
	err := db.Model(enc).
		Where("slug = ?", slug).
		Select()

	if err != nil {
		fmt.Println(err)
		return &models.Encounter{}, err
	}

	fmt.Println("Encounter loaded From DB")
	return enc, nil
}

// DeleteEncounter deletes a single encounter from DB by ID
func DeleteEncounter(db *pg.DB, pk int64) error {

	enc := models.Encounter{ID: pk}

	fmt.Println("Deleting encounter...")

	_, err := db.Model(&enc).Delete()

	return err
}
