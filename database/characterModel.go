package database

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/gosimple/slug"
	"github.com/toferc/rq_web/models"
)

// SaveCharacterModel saves a Character to the DB
func SaveCharacterModel(db *pg.DB, cm *models.CharacterModel) error {

	cm.Character.UpdateCharacter()

	if cm.Slug == "" {
		cm.Slug = slug.Make(fmt.Sprintf("%s-%s", cm.Author.UserName, cm.Character.Name))
	}

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

// UpdateCharacterModel updates a runequest character
func UpdateCharacterModel(db *pg.DB, cm *models.CharacterModel) error {

	cm.Character.UpdateCharacter()

	if cm.Slug == "" {
		cm.Slug = slug.Make(fmt.Sprintf("%s-%s", cm.Author.UserName, cm.Character.Name))
	}

	err := db.Update(cm)
	if err != nil {
		panic(err)
	}
	return err
}

// ListAllCharacterModels queries Character names and add to slice
func ListAllCharacterModels(db *pg.DB) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	_, err := db.Query(&cms, `SELECT * FROM character_models ORDER BY created_at DESC;`)

	if err != nil {
		panic(err)
	}

	return cms, nil
}

// ListCharacterModels queries Character names and add to slice
func ListCharacterModels(db *pg.DB) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	_, err := db.Query(&cms,
		`SELECT * FROM character_models WHERE open = true 
		ORDER BY created_at DESC;`)

	if err != nil {
		panic(err)
	}

	// Print names and PK
	for i, cm := range cms {

		if cm.Character.CreationSteps == nil {
			cm.Character.CreationSteps = map[string]bool{}
			cm.Character.CreationSteps["Complete"] = true
		}

		fmt.Println(i, cm.Character.Name)
	}
	return cms, nil
}

// ListUserCharacterModels queries Character names and add to slice
func ListUserCharacterModels(db *pg.DB, username string) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	_, err := db.Query(&cms,
		`SELECT * FROM character_models WHERE author ->> 'UserName' = ? ORDER BY created_at DESC;`, username)

	if err != nil {
		panic(err)
	}

	return cms, nil
}

func countUserCharacterModels(db *pg.DB, username string) int {

	var count int

	_, err := db.Query(&count,
		`SELECT COUNT(*) FROM character_models WHERE author ->> 'UserName' = ?;`, username)
	if err != nil {
		panic(err)
	}

	return count

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

	if cm.Character.CreationSteps == nil {
		cm.Character.CreationSteps = map[string]bool{}
		cm.Character.CreationSteps["Complete"] = true
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
