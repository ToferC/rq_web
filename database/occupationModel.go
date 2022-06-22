package database

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/gosimple/slug"
	"github.com/toferc/rq_web/models"
	"github.com/toferc/runequest"
)

// SaveOccupationModel saves a Occupation to the DB
func SaveOccupationModel(db *pg.DB, oc *models.OccupationModel) error {

	if oc.Slug == "" {
		oc.Slug = slug.Make(oc.Occupation.Name)
	}

	// Save character in Database
	_, err := db.Model(oc).
		OnConflict("(id) DO UPDATE").
		Set("occupation = ?occupation").
		Insert(oc)
	if err != nil {
		panic(err)
	}
	return err
}

// UpdateOccupationModel updates a Occupation in the database
func UpdateOccupationModel(db *pg.DB, oc *models.OccupationModel) error {

	if oc.Slug == "" {
		oc.Slug = slug.Make(oc.Occupation.Name)
	}

	err := db.Update(oc)
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

	ocMap := map[string]*models.OccupationModel{}

	// Create Map
	for i, oc := range occupations {
		ocMap[runequest.ToSnakeCase(oc.Occupation.Name)] = oc
		fmt.Println(i, oc.Occupation.Name)
	}
	return ocMap, nil
}

// ListOfficialOccupationModels queries Occupation names and add to slice
func ListOfficialOccupationModels(db *pg.DB) (map[string]*models.OccupationModel, error) {
	var occupations []*models.OccupationModel

	_, err := db.Query(&occupations, `SELECT * FROM occupation_models where official = true`)

	if err != nil {
		panic(err)
	}

	ocMap := map[string]*models.OccupationModel{}

	// Create Map
	for i, oc := range occupations {
		ocMap[runequest.ToSnakeCase(oc.Occupation.Name)] = oc
		fmt.Println(i, oc.Occupation.Name)
	}
	return ocMap, nil
}

// LoadOccupationModel loads a single Occupation from the DB by name
func LoadOccupationModel(db *pg.DB, slug string) (*models.OccupationModel, error) {
	// Select user by Primary Key
	occupation := new(models.OccupationModel)
	err := db.Model(occupation).
		Where("Slug = ?", slug).
		Limit(1).
		Select()

	if err != nil {
		return occupation, err
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
