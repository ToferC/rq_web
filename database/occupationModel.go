package database

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/toferc/rq_web/models"
	"github.com/toferc/runequest"
)

// SaveOccupationModel saves a Occupation to the DB
func SaveOccupationModel(db *pg.DB, hl *models.OccupationModel) error {

	// Save character in Database
	_, err := db.Model(hl).
		OnConflict("(id) DO UPDATE").
		Set("occupation = ?occupation").
		Insert(hl)
	if err != nil {
		panic(err)
	}
	return err
}

// UpdateOccupationModel updates a Occupation in the database
func UpdateOccupationModel(db *pg.DB, hl *models.OccupationModel) error {

	err := db.Update(hl)
	if err != nil {
		panic(err)
	}
	return err
}

// ListOccupationModels queries Occupation names and add to slice
func ListOccupationModels(db *pg.DB) (map[string]*models.OccupationModel, error) {
	var occupations []*models.OccupationModel

	_, err := db.Query(&occupations, `SELECT * FROM occupation_models`)

	if err != nil {
		panic(err)
	}

	hlMap := map[string]*models.OccupationModel{}

	// Create Map
	for i, hl := range occupations {
		hlMap[runequest.ToSnakeCase(hl.Occupation.Name)] = hl
		fmt.Println(i, hl.Occupation.Name)
	}
	return hlMap, nil
}

// LoadOccupationModel loads a single Occupation from the DB by name
func LoadOccupationModel(db *pg.DB, name string) (*models.OccupationModel, error) {
	// Select user by Primary Key
	occupation := new(models.OccupationModel)
	err := db.Model(occupation).
		Where("Name = ?", name).
		Limit(1).
		Select()

	if err != nil {
		panic(err)
	}

	fmt.Println("Occupation loaded From DB")
	return occupation, nil
}

// PKLoadOccupationModel loads a single Occupation from the DB by pk
func PKLoadOccupationModel(db *pg.DB, pk int64) (*models.OccupationModel, error) {
	// Select user by Primary Key
	occupation := &models.OccupationModel{ID: pk}
	err := db.Select(occupation)

	if err != nil {
		return &models.OccupationModel{}, err
	}

	fmt.Println("Occupation loaded From DB")
	return occupation, nil
}

// DeleteOccupationModel deletes a single Occupation from DB by ID
func DeleteOccupationModel(db *pg.DB, pk int64) error {

	pow := models.OccupationModel{ID: pk}

	fmt.Println("Deleting Occupation...")

	err := db.Delete(&pow)

	return err
}
