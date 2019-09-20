package database

import (
	"fmt"
	"log"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/toferc/rq_web/models"
)

// InitDB initializes the DB Schema
func InitDB(db *pg.DB) error {
	err := createSchema(db)
	if err != nil {
		panic(err)
	}
	return err
}

// CreateIndex creates an index for TSVs on character models
func CreateIndex(db *pg.DB) error {
	fmt.Println("Creating TSV Index")
	_, err := db.Exec("CREATE INDEX IF NOT EXISTS idx_fts_doc_vec ON character_models USING gin(tsv);")
	if err != nil {
		log.Println(err)
	}
	return err
}

// CreateTSVColumn adds tsv on character models
func CreateTSVColumn(db *pg.DB) error {
	fmt.Println("Creating TSV Column on character_models")
	_, err := db.Exec("ALTER TABLE character_models ADD COLUMN IF NOT EXISTS tsv tsvector;")
	if err != nil {
		log.Println(err)
	}
	return err
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{
		(*models.CharacterModel)(nil),
		(*models.HomelandModel)(nil),
		(*models.OccupationModel)(nil),
		(*models.CultModel)(nil),
		(*models.Image)(nil),
		(*models.Faction)(nil),
		(*models.Encounter)(nil),
		(*models.Campaign)(nil),
		(*models.User)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func counter(ar []string) map[string]int {
	m := map[string]int{}

	for _, a := range ar {
		m[a]++
	}

	return m
}
