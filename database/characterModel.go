package database

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/models"
)

// SaveCharacterModel saves a Character to the DB
func SaveCharacterModel(db *pg.DB, cm *models.CharacterModel) error {

	oneroll.UpdateCost(cm.Character)

	// Save character in Database
	_, err := db.Model(cm).
		OnConflict("(id) DO UPDATE").
		Set("character = ?character").
		Insert(cm)
	if err != nil {
		panic(err)
	}
	return err
}

func UpdateCharacterModel(db *pg.DB, cm *models.CharacterModel) error {

	oneroll.UpdateCost(cm.Character)

	err := db.Update(cm)
	if err != nil {
		panic(err)
	}
	return err
}

// ListCharacterModels queries Character names and add to slice
func ListCharacterModels(db *pg.DB) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	_, err := db.Query(&cms, `SELECT * FROM character_models WHERE open = true`)

	if err != nil {
		panic(err)
	}

	// Print names and PK
	for i, n := range cms {
		fmt.Println(i, n.Character.Name)
	}
	return cms, nil
}

// ListUserCharacterModels queries Character names and add to slice
func ListUserCharacterModels(db *pg.DB, username string) ([]*models.CharacterModel, error) {
	var temp []*models.CharacterModel
	var cms []*models.CharacterModel

	_, err := db.Query(&temp, `SELECT * FROM character_models`)

	if err != nil {
		panic(err)
	}

	for _, t := range temp {
		if t.Author.UserName == username {
			cms = append(cms, t)
		}
	}

	// Print names and PK
	for i, n := range cms {
		fmt.Println(i, n.Character.Name)
	}
	return cms, nil
}

// PKLoadCharacterModel loads a single character from the DB by pk
func PKLoadCharacterModel(db *pg.DB, pk int64) (*models.CharacterModel, error) {
	// Select user by Primary Key
	cm := &models.CharacterModel{ID: pk}
	err := db.Select(cm)

	if err != nil {
		fmt.Println(err)
		return &models.CharacterModel{}, err
	}

	fmt.Println("Character loaded From DB")
	return cm, nil
}

// DeleteCharacterModel deletes a single character from DB by ID
func DeleteCharacterModel(db *pg.DB, pk int64) error {

	cm := models.CharacterModel{ID: pk}

	fmt.Println("Deleting character...")

	err := db.Delete(&cm)

	return err
}
