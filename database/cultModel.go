package database

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/gosimple/slug"
	"github.com/toferc/rq_web/models"
	"github.com/toferc/runequest"
)

// SaveCultModel saves a Cult to the DB
func SaveCultModel(db *pg.DB, cl *models.CultModel) error {

	if cl.Slug == "" {
		cl.Slug = slug.Make(cl.Cult.Name)
	}

	// Save character in Database
	_, err := db.Model(cl).
		OnConflict("(id) DO UPDATE").
		Set("cult = ?cult").
		Insert(cl)
	if err != nil {
		panic(err)
	}
	return err
}

// UpdateCultModel updates a Cult in the database
func UpdateCultModel(db *pg.DB, cl *models.CultModel) error {

	if cl.Slug == "" {
		cl.Slug = slug.Make(cl.Cult.Name)
	}

	err := db.Update(cl)
	if err != nil {
		panic(err)
	}
	return err
}

// ListCultModels queries Cult names and add to slice
func ListCultModels(db *pg.DB) (map[string]*models.CultModel, error) {
	var cults []*models.CultModel

	_, err := db.Query(&cults, `SELECT * FROM cult_models`)

	if err != nil {
		panic(err)
	}

	clMap := map[string]*models.CultModel{}

	// Create Map
	for _, cl := range cults {
		clMap[runequest.ToSnakeCase(cl.Cult.Name)] = cl
	}
	return clMap, nil
}

// ListOfficialCultModels queries Cult names and add to slice
func ListOfficialCultModels(db *pg.DB) (map[string]*models.CultModel, error) {
	var cults []*models.CultModel

	_, err := db.Query(&cults, `SELECT * FROM cult_models WHERE official = true`)

	if err != nil {
		panic(err)
	}

	clMap := map[string]*models.CultModel{}

	// Create Map
	for _, cl := range cults {
		clMap[runequest.ToSnakeCase(cl.Cult.Name)] = cl
	}
	return clMap, nil
}

// LoadCultModel loads a single Cult from the DB by name
func LoadCultModel(db *pg.DB, slug string) (*models.CultModel, error) {
	// Select user by Primary Key
	cult := new(models.CultModel)
	err := db.Model(cult).
		Where("Slug = ?", slug).
		Limit(1).
		Select()

	if err != nil {
		return cult, err
	}

	fmt.Println("Cult loaded From DB")
	return cult, nil
}

// PKLoadCultModel loads a single Cult from the DB by pk
func PKLoadCultModel(db *pg.DB, pk int64) (*models.CultModel, error) {
	// Select user by Primary Key
	cult := &models.CultModel{ID: pk}
	err := db.Select(cult)

	if err != nil {
		return &models.CultModel{}, err
	}

	fmt.Println("Cult loaded From DB")
	return cult, nil
}

// DeleteCultModel deletes a single Cult from DB by ID
func DeleteCultModel(db *pg.DB, pk int64) error {

	pow := models.CultModel{ID: pk}

	fmt.Println("Deleting Cult...")

	err := db.Delete(&pow)

	return err
}
