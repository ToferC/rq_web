package model

// Action describes an action that a service conducts
type Action struct {
	ID           int64
	Name         string
	Description  string
	Service      Service
	TargetUser   *User
	TargetModule *Module
	TargetField  map[string]*Field
}
