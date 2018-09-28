package model

import (
	"time"
)

// Profile describes a single user profile across OneGC
type Profile struct {
	ID              int64
	Name            string
	DOB             string
	TitleEn         string
	TitleFr         string
	Location        *Point
	Organization    *Organization
	Modules         map[string]*Module
	BasePermissions map[string]*Permission
}

// Module defines a special profile datastore of fields
type Module struct {
	ID          int64
	Name        string
	Description string
	Fields      []*Field
}

// Field describes a Profile field within a Module
type Field struct {
	ID          int64
	Name        string
	Data        string
	Permissions []*Permission
	AccessLogs  []*AccessLog
}

// AccessLog describes and logs Service access of Profile Fields
type AccessLog struct {
	ID          int64
	Time        time.Time
	Service     Service
	Action      Action
	Description string
}

// Permission tracks User's granting Services access to fields
type Permission struct {
	ID      int64
	Service Service
	Actions []Action
	Time    *time.Time
}
