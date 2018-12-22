package database

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/toferc/rq_web/models"
)

// SaveCreatureModel saves a Creature to the DB
func SaveCreatureModel(db *pg.DB, cm *models.CreatureModel) error {

	cm.Creature.UpdateCreature()

	// Save creature in Database
	_, err := db.Model(cm).
		OnConflict("(id) DO UPDATE").
		Set("creature = ?creature").
		Insert(cm)
	if err != nil {
		panic(err)
	}
	return err
}

// UpdateCreatureModel updates a runequest creature
func UpdateCreatureModel(db *pg.DB, cm *models.CreatureModel) error {

	cm.Creature.UpdateCreature()

	err := db.Update(cm)
	if err != nil {
		panic(err)
	}
	return err
}

// ListCreatureModels queries Creature names and add to slice
func ListCreatureModels(db *pg.DB) ([]*models.CreatureModel, error) {
	var cms []*models.CreatureModel

	_, err := db.Query(&cms, `SELECT * FROM creature_models WHERE open = true`)

	if err != nil {
		panic(err)
	}

	// Print names and PK
	for i, cm := range cms {

		fmt.Println(i, cm.Creature.Name)
	}
	return cms, nil
}

// ListUserCreatureModels queries Creature names and add to slice
func ListUserCreatureModels(db *pg.DB, username string) ([]*models.CreatureModel, error) {
	var temp []*models.CreatureModel
	var cms []*models.CreatureModel

	_, err := db.Query(&temp, `SELECT * FROM creature_models`)

	if err != nil {
		panic(err)
	}

	for _, t := range temp {
		if t.Author.UserName == username {
			cms = append(cms, t)
		}
	}

	// Print names and PK
	for i, cm := range cms {

		fmt.Println(i, cm.Creature.Name)
	}
	return cms, nil
}

// PKLoadCreatureModel loads a single creature from the DB by pk
func PKLoadCreatureModel(db *pg.DB, pk int64) (*models.CreatureModel, error) {
	// Select user by Primary Key
	cm := &models.CreatureModel{ID: pk}
	err := db.Select(cm)

	if err != nil {
		fmt.Println(err)
		return &models.CreatureModel{}, err
	}

	fmt.Println("Creature loaded From DB")
	return cm, nil
}

// DeleteCreatureModel deletes a single creature from DB by ID
func DeleteCreatureModel(db *pg.DB, pk int64) error {

	cm := models.CreatureModel{ID: pk}

	fmt.Println("Deleting creature...")

	err := db.Delete(&cm)

	return err
}
