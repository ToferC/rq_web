package database

import (
	"fmt"
	"log"

	"github.com/go-pg/pg"
	"github.com/toferc/rq_web/models"
)

// SaveFaction saves a Faction to the DB
func SaveFaction(db *pg.DB, fac *models.Faction) error {

	// Save faction in Database
	_, err := db.Model(fac).
		OnConflict("(id) DO UPDATE").
		Set("name = ?name").
		Insert(fac)
	if err != nil {
		panic(err)
	}
	return err
}

// UpdateFaction updates a runequest faction
func UpdateFaction(db *pg.DB, fac *models.Faction) error {

	err := db.Update(fac)
	if err != nil {
		panic(err)
	}
	return err
}

// ListFactions queries Faction names and add to slice
func ListFactions(db *pg.DB) ([]*models.Faction, error) {
	var facs []*models.Faction

	_, err := db.Query(&facs, `SELECT * FROM factions`)

	if err != nil {
		panic(err)
	}

	return facs, nil
}

// ListUserFactions queries Faction names and add to slice
func ListUserFactions(db *pg.DB, username string) ([]*models.Faction, error) {
	var facs []*models.Faction

	_, err := db.Query(&facs, `SELECT * FROM factions WHERE author ->> 'UserName' = ?`, username)

	if err != nil {
		panic(err)
	}

	return facs, nil
}

// PKLoadFaction loads a single faction from the DB by pk
func PKLoadFaction(db *pg.DB, pk int64) (*models.Faction, error) {
	// Select user by Primary Key
	fac := &models.Faction{ID: pk}
	err := db.Select(fac)

	if err != nil {
		fmt.Println(err)
		return &models.Faction{}, err
	}

	fmt.Println("Faction loaded From DB")
	return fac, nil
}

// SlugLoadFaction loads a single faction from the DB by pk
func SlugLoadFaction(db *pg.DB, slug string) (*models.Faction, error) {
	// Select user by Primary Key
	fac := &models.Faction{}
	err := db.Model(fac).
		Where("slug = ?", slug).
		Select()

	if err != nil {
		fmt.Println(err)
		return &models.Faction{}, err
	}

	fmt.Println("Faction loaded From DB")
	return fac, nil
}

// LoadFactionCharacterModels queries Character names and add to slice
func LoadFactionCharacterModels(db *pg.DB, slugs []string) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	_, err := db.Query(&cms, `SELECT * FROM character_models WHERE slug IN (?)`, pg.In(slugs))

	if err != nil {
		log.Panic(err)
	}

	c := counter(slugs)

	for k, v := range c {
		if v > 1 {
			for i := 1; i < v; i++ {
				for _, cm := range cms {
					if cm.Slug == k {
						cms = append(cms, cm)
						break
					}
				}
			}
		}
	}

	return cms, nil
}

func counter(ar []string) map[string]int {
	m := map[string]int{}

	for _, a := range ar {
		m[a]++
	}

	return m
}

// DeleteFaction deletes a single faction from DB by ID
func DeleteFaction(db *pg.DB, pk int64) error {

	fac := models.Faction{ID: pk}

	fmt.Println("Deleting faction...")

	err := db.Delete(&fac)

	return err
}
