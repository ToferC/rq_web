package database

import (
	"fmt"
	"log"

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
		log.Println(err)
	}

	setTSVSearch := fmt.Sprintf(`
	UPDATE character_models SET search = 
	setweight(to_tsvector(coalesce(character ->> 'Name')), 'A') ||
	setweight(to_tsvector(coalesce(author ->> 'UserName')), 'B') ||
	setweight(to_tsvector(coalesce(character ->> 'Description')), 'C') ||

	setweight(to_tsvector(coalesce(character#> '{Homeland, Name}')), 'D') ||
	setweight(to_tsvector(coalesce(character#> '{Occupation, Name}')), 'D') ||
	setweight(to_tsvector(coalesce(character#> '{Cult, Name}')), 'D') ||

	WHERE
	id = %d;
	`, cm.ID)

	db.ExecOne(setTSVSearch)

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
		log.Println(err)
	}

	setTSVSearch := fmt.Sprintf(`
	UPDATE character_models SET search = 
	setweight(to_tsvector(coalesce(character ->> 'Name')), 'A') ||
	setweight(to_tsvector(coalesce(author ->> 'UserName')), 'B') ||
	setweight(to_tsvector(coalesce(character ->> 'Description')), 'C') ||

	setweight(to_tsvector(coalesce(character#> '{Homeland, Name}')), 'D') ||
	setweight(to_tsvector(coalesce(character#> '{Occupation, Name}')), 'D') ||
	setweight(to_tsvector(coalesce(character#> '{Cult, Name}')), 'D') ||

	WHERE
	id = %d;
	`, cm.ID)

	db.ExecOne(setTSVSearch)

	return err
}

// ListAllCharacterModels queries Character names and add to slice
func ListAllCharacterModels(db *pg.DB) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	_, err := db.Query(&cms, `SELECT * FROM character_models ORDER BY created_at DESC;`)

	if err != nil {
		log.Println(err)
	}

	return cms, nil
}

// SearchCharacterModels queries Character names and add to slice
func SearchCharacterModels(db *pg.DB, q string) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	_, err := db.Query(&cms, `
				SELECT *,
				ts_rank_cd(search, q) AS RANK
				FROM character_models, plainto_tsquery(?) q
				WHERE
				search @@ q AND open = 'true'
				ORDER BY rank DESC
				LIMIT 25;`, q)

	if err != nil {
		log.Println(err)
	}

	return cms, nil
}

// ListCharacterModels queries open Character names and add to slice
func ListCharacterModels(db *pg.DB) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	_, err := db.Query(&cms, `SELECT * FROM character_models ORDER BY created_at DESC
						WHERE open = true;`)

	if err != nil {
		log.Println(err)
	}

	return cms, nil
}

// PaginateCharacterModels queries open Character names and add to slice
func PaginateCharacterModels(db *pg.DB, limit, offset int) ([]*models.CharacterModel, error) {

	var cms []*models.CharacterModel

	err := db.Model(&cms).
		Limit(limit).
		Offset(offset).
		Where("open = true").
		Order("created_at DESC").
		Select()

	if err != nil {
		log.Println(err)
	}

	return cms, nil
}

// ListCraftedCharacterModels queries open Character names and add to slice
func ListCraftedCharacterModels(db *pg.DB) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	_, err := db.Query(&cms,
		`SELECT * FROM character_models WHERE open = true AND random = false
		UNION
		SELECT * FROM character_models WHERE open = true AND random IS NULL

		ORDER BY created_at DESC;`)

	if err != nil {
		log.Println(err)
	}

	return cms, nil
}

// ListRandomCharacterModels queries open Character names and add to slice
func ListRandomCharacterModels(db *pg.DB) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	_, err := db.Query(&cms,
		`SELECT * FROM character_models WHERE open = true AND random = true
		ORDER BY created_at DESC;`)

	if err != nil {
		log.Println(err)
	}

	return cms, nil
}

// ListUserCharacterModels queries Character names and add to slice
func ListUserCharacterModels(db *pg.DB, username string) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	_, err := db.Query(&cms,
		`SELECT * FROM character_models WHERE author ->> 'UserName' = ? ORDER BY created_at DESC;`, username)

	if err != nil {
		log.Println(err)
	}

	return cms, nil
}

func countUserCharacterModels(db *pg.DB, username string) int {

	var count int

	_, err := db.Query(&count,
		`SELECT COUNT(*) FROM character_models WHERE author ->> 'UserName' = ?;`, username)
	if err != nil {
		log.Println(err)
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
