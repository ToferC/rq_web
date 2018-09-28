package database

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/models"
)

// SavePowerModel saves a PowerModel to the DB
func SavePowerModel(db *pg.DB, pm *models.PowerModel) error {

	// Save character in Database
	_, err := db.Model(pm).
		OnConflict("(id) DO UPDATE").
		Set("power = ?power").
		Insert(pm)
	if err != nil {
		panic(err)
	}
	return err
}

func UpdatePowerModel(db *pg.DB, pm *models.PowerModel) error {

	err := db.Update(pm)
	if err != nil {
		panic(err)
	}
	return err
}

// ListPowerModels queries PowerModel names and add to slice
func ListPowerModels(db *pg.DB) (map[string]models.PowerModel, error) {
	var pows []models.PowerModel

	_, err := db.Query(&pows, `SELECT * FROM power_models`)

	if err != nil {
		panic(err)
	}

	powMap := map[string]models.PowerModel{}

	// Create Map
	for i, p := range pows {
		powMap[oneroll.ToSnakeCase(p.Power.Name)] = p
		fmt.Println(i, p.Power.Name)
	}
	return powMap, nil
}

// PKLoadPowerModel loads a single PowerModel from the DB by pk
func PKLoadPowerModel(db *pg.DB, pk int64) (*models.PowerModel, error) {
	// Select user by Primary Key
	pow := &models.PowerModel{ID: pk}
	err := db.Select(pow)

	if err != nil {
		return &models.PowerModel{}, err
	}

	fmt.Println("PowerModel loaded From DB")
	return pow, nil
}

// DeletePowerModel deletes a single PowerModel from DB by ID
func DeletePowerModel(db *pg.DB, pk int64) error {

	pow := models.PowerModel{ID: pk}

	fmt.Println("Deleting PowerModel...")

	err := db.Delete(&pow)

	return err
}
