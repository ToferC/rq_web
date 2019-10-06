package models

import "time"

// Note represents a user's game notes
type Note struct {
	ID               int64 `schema:"-"`
	CharacterModelID int64 `schema:"-"`
	AuthorID         int64 `schema:"-"`
	AuthorUserName   string
	Year             int
	Season           string
	Week             string
	Title            string
	Body             string
	Likes            int       `schema:"-"`
	PublishedOn      time.Time `schema:"-"`
	Open             bool
	Tags             []string
	Slug             string `schema:"-"`
}

// Seasons represents the model of time in Glorantha
var Seasons = []string{
	"Sea", "Fire", "Earth", "Darkness", "Storm", "Sacred Time",
}

// Weeks are the weeks of the Gloranthan calendar
var Weeks = []string{
	"Disorder", "Harmony", "Death", "Fertility", "Stasis",
	"Movement", "Illusion", "Truth",
}
