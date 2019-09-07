package database

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/gosimple/slug"
	"github.com/toferc/rq_web/models"
	"github.com/toferc/runequest"
)

// SaveHomelandModel saves a Homeland to the DB
func SaveHomelandModel(db *pg.DB, hl *models.HomelandModel) error {

	if hl.Slug == "" {
		hl.Slug = slug.Make(fmt.Sprintf("%s", hl.Homeland.Name))
	}

	// Save character in Database
	_, err := db.Model(hl).
		OnConflict("(id) DO UPDATE").
		Set("homeland = ?homeland").
		Insert(hl)
	if err != nil {
		panic(err)
	}
	return err
}

// UpdateHomelandModel updates a Homeland in the database
func UpdateHomelandModel(db *pg.DB, hl *models.HomelandModel) error {

	if hl.Slug == "" {
		hl.Slug = slug.Make(fmt.Sprintf("%s", hl.Homeland.Name))
	}

	err := db.Update(hl)
	if err != nil {
		panic(err)
	}
	return err
}

// ListHomelandModels queries Homeland names and add to slice
func ListHomelandModels(db *pg.DB) (map[string]*models.HomelandModel, error) {
	var homelands []*models.HomelandModel

	_, err := db.Query(&homelands, `SELECT * FROM homeland_models`)

	if err != nil {
		panic(err)
	}

	hlMap := map[string]*models.HomelandModel{}

	// Create Map
	for i, hl := range homelands {
		hlMap[runequest.ToSnakeCase(hl.Homeland.Name)] = hl
		fmt.Println(i, hl.Homeland.Name)
	}
	return hlMap, nil
}

// ListOfficialHomelandModels queries Homeland names and add to slice
func ListOfficialHomelandModels(db *pg.DB) (map[string]*models.HomelandModel, error) {
	var homelands []*models.HomelandModel

	_, err := db.Query(&homelands, `SELECT * FROM homeland_models WHERE official = true`)

	if err != nil {
		panic(err)
	}

	hlMap := map[string]*models.HomelandModel{}

	// Create Map
	for i, hl := range homelands {
		hlMap[runequest.ToSnakeCase(hl.Homeland.Name)] = hl
		fmt.Println(i, hl.Homeland.Name)
	}
	return hlMap, nil
}

// LoadHomelandModel loads a single Homeland from the DB by name
func LoadHomelandModel(db *pg.DB, slug string) (*models.HomelandModel, error) {
	// Select user by Primary Key
	homeland := new(models.HomelandModel)
	err := db.Model(homeland).
		Where("Slug = ?", slug).
		Limit(1).
		Select()

	if err != nil {
		return homeland, err
	}

	fmt.Println("Homeland loaded From DB")
	return homeland, nil
}

// PKLoadHomelandModel loads a single Homeland from the DB by pk
func PKLoadHomelandModel(db *pg.DB, pk int64) (*models.HomelandModel, error) {
	// Select user by Primary Key
	homeland := &models.HomelandModel{ID: pk}
	err := db.Select(homeland)

	if err != nil {
		return &models.HomelandModel{}, err
	}

	fmt.Println("Homeland loaded From DB")
	return homeland, nil
}

// DeleteHomelandModel deletes a single Homeland from DB by ID
func DeleteHomelandModel(db *pg.DB, pk int64) error {

	pow := models.HomelandModel{ID: pk}

	fmt.Println("Deleting Homeland...")

	err := db.Delete(&pow)

	return err
}
