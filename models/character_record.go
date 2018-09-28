package models

import (
	"time"

	"github.com/toferc/runequest"
)

type CharacterModel struct {
	ID        int64
	Author    *User
	Character *runequest.Character
	Open      bool
	Likes     int
	Image     *Image
	Slug      string
	CreatedAt time.Time `sql:"default:now()"`
	UpdatedAt time.Time
}

type Image struct {
	Id   int
	Path string
}
