package database

import (
	"fmt"
	"log"

	"github.com/go-pg/pg"
	"github.com/toferc/rq_web/models"
)

// SaveNote saves a Note to the DB
func SaveNote(db *pg.DB, nt *models.Note) error {

	// Save note in Database
	_, err := db.Model(nt).
		OnConflict("(id) DO UPDATE").
		Set("title = ?title").
		Insert(nt)
	if err != nil {
		panic(err)
	}
	return err
}

// UpdateNote updates a runequest note
func UpdateNote(db *pg.DB, nt *models.Note) error {

	err := db.Update(nt)
	if err != nil {
		panic(err)
	}
	return err
}

// ListNotes queries Note names and add to slice
func ListNotes(db *pg.DB, cmID int64) ([]*models.Note, error) {
	var nts []*models.Note

	_, err := db.Query(&nts,
		`SELECT * FROM notes
		 WHERE character_model_id = ?
		 ORDER BY year DESC,
		 CASE WHEN season = 'Sea' THEN 1
			WHEN season = 'Fire' THEN 2 
			WHEN season = 'Earth' THEN 3 
			WHEN season = 'Darkness' THEN 4 
			WHEN season = 'Storm' THEN 5
			WHEN season = 'Sacred Time' THEN 6
		 END DESC,
		 CASE WHEN week = 'Disorder' THEN 1
		 	WHEN week = 'Harmony' THEN 2
		 	WHEN week = 'Death' THEN 3
		 	WHEN week = 'Fertility' THEN 4
		 	WHEN week = 'Stasis' THEN 5
		 	WHEN week = 'Movement' THEN 6
		 	WHEN week = 'Illusion' THEN 7
		 	WHEN week = 'Truth' THEN 8
		END DESC;`, cmID)

	if err != nil {
		log.Println(err)
	}

	return nts, err
}

func countNotes(db *pg.DB, id int64) int {

	var count int

	_, err := db.Query(&count,
		`SELECT COUNT(*) FROM notes WHERE character_model_id = ?;`, id)
	if err != nil {
		log.Println(err)
	}

	return count

}

// PKLoadNote loads a single note from the DB by pk
func PKLoadNote(db *pg.DB, pk int64) (*models.Note, error) {
	// Select user by Primary Key
	nt := &models.Note{ID: pk}
	err := db.Select(nt)

	if err != nil {
		fmt.Println(err)
		return &models.Note{}, err
	}

	fmt.Println("Note loaded From DB")
	return nt, nil
}

// SlugLoadNote loads a single note from the DB by pk
func SlugLoadNote(db *pg.DB, slug string) (*models.Note, error) {
	// Select user by Primary Key
	nt := &models.Note{}
	err := db.Model(nt).
		Where("slug = ?", slug).
		Select()

	if err != nil {
		fmt.Println(err)
		return &models.Note{}, err
	}

	fmt.Println("Note loaded From DB")
	return nt, nil
}

// DeleteNote deletes a single note from DB by ID
func DeleteNote(db *pg.DB, pk int64) error {

	nt := models.Note{ID: pk}

	fmt.Println("Deleting note...")

	err := db.Delete(&nt)

	return err
}
