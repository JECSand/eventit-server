package models

import (
	"github.com/JECSand/eventit-server/domains/shared/enums"
	"time"
)

// User is a root struct that is used to store the json encoded data for/from a mongodb user doc.
type User struct {
	Id        string     `json:"id,omitempty"`
	Username  string     `json:"username,omitempty"`
	Password  string     `json:"password,omitempty"`
	FirstName string     `json:"firstname,omitempty"`
	LastName  string     `json:"lastname,omitempty"`
	Email     string     `json:"email,omitempty"`
	Role      enums.Role `json:"role,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	DeletedAt time.Time  `json:"deleted_at,omitempty"`
}
