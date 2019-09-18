package models

import (
	"time"

	"github.com/toferc/runequest"
)

// CharacterModel represents the Web struct for a runequest.Character
// Must add tsv column tsvector type to database for search
type CharacterModel struct {
	tableName struct{} `pg:",discard_unknown_columns"`

	ID        int64
	Author    *User
	Random    bool
	Character *runequest.Character
	Open      bool
	Likes     int
	LikeData  map[string]*Like
	Image     *Image
	Slug      string
	CreatedAt time.Time `sql:"default:now()"`
	UpdatedAt time.Time
}

// Image  is the image and path for an Image
type Image struct {
	ID   int
	Path string
}

// Like represents a user appreciation of an object
type Like struct {
	UserName  string
	CreatedAt time.Time `sql:"default:now()"`
}
