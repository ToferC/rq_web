package models

import (
	"time"

	"github.com/toferc/runequest"
)

// CreatureModel represents the Web struct for a runequest.creature
type CreatureModel struct {
	ID        int64
	Author    *User
	Creature  *runequest.Creature
	Open      bool
	Likes     int
	Image     *Image
	Slug      string
	CreatedAt time.Time `sql:"default:now()"`
	UpdatedAt time.Time
}
