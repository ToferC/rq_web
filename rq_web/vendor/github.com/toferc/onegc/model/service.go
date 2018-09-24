package model

// Service describes a service offering in OneGC
type Service struct {
	ID          int64
	Name        string
	Description string
	Image       string
	Slug        string
	Actions     map[string]*Action
}
