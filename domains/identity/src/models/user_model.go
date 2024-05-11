package models

import (
	"errors"
	"github.com/JECSand/eventit-server/domains/shared/enums"
	"golang.org/x/crypto/bcrypt"
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

// HashPassword hashes a user password and associates it with the user struct
func (g *User) HashPassword() error {
	if len(g.Password) != 0 {
		password := []byte(g.Password)
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		g.Password = string(hashedPassword)
		return nil
	}
	return errors.New("no password set to hash in user model")
}

// UsersPage Multiple Users in a paginated response
type UsersPage struct {
	TotalCount int64   `json:"total_count"`
	TotalPages int64   `json:"total_pages"`
	Page       int64   `json:"page"`
	Size       int64   `json:"size"`
	HasMore    bool    `json:"has_more"`
	Users      []*User `json:"users"`
}

/*
// ToProto convert UsersRes to proto
func (p *UsersRes) ToProto() []*usersService.User {
	uList := make([]*usersService.User, 0, len(p.Users))
	for _, u := range p.Users {
		u.Password = ""
		uList = append(uList, u.ToProto())
	}
	return uList
}
*/
