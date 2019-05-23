package models

import "time"

// Encounter represents a combat situation
type Encounter struct {
	ID          int64
	Author      *User
	Name        string
	Description string
	Factions    []*Faction
	Outcome     string

	Open      bool
	Likes     int
	Image     *Image
	Slug      string
	CreatedAt time.Time `sql:"default:now()"`
	UpdatedAt time.Time
}

// Faction represents a side in an Encounter
type Faction struct {
	ID                  int64
	Author              *User
	Name                string
	Description         string
	CharacterModelSlugs []string
	Open                bool
	Likes               int
	Image               *Image
	Slug                string
	CreatedAt           time.Time `sql:"default:now()"`
	UpdatedAt           time.Time
}
