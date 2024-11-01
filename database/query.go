package database

import (
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/toferc/rq_web/models"
)

// QueryArgs is the base for a filtered database query
type QueryArgs struct {
	UserName   string
	Homeland   string
	Occupation string
	Cult       string
}

// SearchCharacterModels queries Character names and add to slice
func SearchCharacterModels(db *pg.DB, q string) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	_, err := db.Query(&cms, `
				SELECT *,
				ts_rank_cd(tsv, q) AS RANK
				FROM character_models, plainto_tsquery(?) q
				WHERE
				tsv @@ q AND open = 'true'
				ORDER BY rank DESC
				;`, q)

	if err != nil {
		log.Println(err)
	}

	return cms, nil
}

// GetFilteredCharacterModels returns characters matching the query
func (q *QueryArgs) GetFilteredCharacterModels(db *pg.DB) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	if q.Homeland == "" {
		q.Homeland = "%"
	}

	if q.Occupation == "" {
		q.Occupation = "%"
	}

	if q.Cult == "" {
		q.Cult = "%"
	}

	_, err := db.Query(&cms, `SELECT * FROM character_models WHERE character -> 'Homeland' ->> 'Name' LIKE ?
		AND character -> 'Occupation' ->> 'Name' LIKE ?
		AND character -> 'Cult' ->> 'Name' LIKE ? 
		AND open = 'true'`, q.Homeland, q.Occupation, q.Cult)

	if err != nil {
		panic(err)
	}

	return cms, nil
}

// GetUserFilteredCharacterModels returns characters matching the query
func (q *QueryArgs) GetUserFilteredCharacterModels(db *pg.DB) ([]*models.CharacterModel, error) {
	var cms []*models.CharacterModel

	if q.Homeland == "" {
		q.Homeland = "%"
	}

	if q.Occupation == "" {
		q.Occupation = "%"
	}

	if q.Cult == "" {
		q.Cult = "%"
	}

	_, err := db.Query(&cms, `SELECT * FROM character_models WHERE character -> 'Homeland' ->> 'Name' LIKE ?
		AND character -> 'Occupation' ->> 'Name' LIKE ?
		AND character -> 'Cult' ->> 'Name' LIKE ? 
		AND author ->> 'UserName' = ?`, q.Homeland, q.Occupation, q.Cult, q.UserName)

	if err != nil {
		panic(err)
	}

	return cms, nil
}
