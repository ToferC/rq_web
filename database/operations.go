package database

import (
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

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{
		(*models.CharacterModel)(nil),
		(*models.CreatureModel)(nil),
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
