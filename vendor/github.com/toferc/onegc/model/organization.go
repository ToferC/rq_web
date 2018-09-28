package model

import (
	"fmt"
)

// Organization describes an organization that provides services
type Organization struct {
	ID          int64
	Name        string
	URL         string
	Description string
	Image       string
	Geo         *Point
	Certificate string
	Secret      string
	Slug        string
	Services    map[string]*Service
}

// Point represents a geocode
type Point struct {
	X, Y float32
}

func (o *Organization) String() string {

	text := fmt.Sprintf("%d: %s \nDescription: %s", o.ID, o.Name, o.Description)

	return text
}
