package database

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/toferc/oneroll"
	"github.com/toferc/ore_web_roller/models"
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
		(*oneroll.Character)(nil),
		(*oneroll.Power)(nil),

		(*models.CharacterModel)(nil),
		(*models.Image)(nil),
		(*models.PowerModel)(nil),
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
