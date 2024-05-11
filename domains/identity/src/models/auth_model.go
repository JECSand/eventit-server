package models

import (
	"errors"
	"github.com/JECSand/eventit-server/domains/shared/auth"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Blacklist is a root struct that is used to store the json encoded data for/from a mongodb blacklist doc.
type Blacklist struct {
	Id        string    `json:"id,omitempty"`
	AuthToken string    `json:"auth_token,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type Credentials struct {
	Email           string `json:"email,omitempty"`
	Password        string `json:"password,omitempty"`
	CurrentPassword string `json:"current_password,omitempty"`
}

type Auth struct {
	User      *User         `json:"user,omitempty"`
	AuthToken string        `json:"auth_token,omitempty"`
	Session   *auth.Session `json:"session,omitempty"`
	CreatedAt time.Time     `json:"created_at,omitempty"`
}

// authenticate compares an input password with the hashed password stored in the User model
func (r *Auth) authenticate(user *User, checkPassword string) error {
	if len(user.Password) != 0 {
		password := []byte(user.Password)
		cPassword := []byte(checkPassword)
		return bcrypt.CompareHashAndPassword(password, cPassword)
	}
	return errors.New("no password set to hash in user model")
}

// Invalidate compares an input password with the hashed password stored in the User model
func (r *Auth) Invalidate() {
	r.AuthToken = ""
	r.User = nil
	r.Session = nil
	return
}

// LoadSession compares an input password with the hashed password stored in the User model
func (r *Auth) LoadSession() (err error) {
	if r.AuthToken == "" {
		err = errors.New("token is empty")
		return
	}
	r.Session, err = auth.LoadSession(r.AuthToken)
	return
}

// NewSession compares an input password with the hashed password stored in the User model
func (r *Auth) NewSession() (err error) {
	if r.User == nil {
		err = errors.New("user is nil")
		return
	}
	r.Session = auth.NewSession(r.User.Id, r.User.Role)
	return
}

// NewToken compares an input password with the hashed password stored in the User model
func (r *Auth) NewToken() (err error) {
	if r.Session == nil {
		err = errors.New("session is nil")
		return
	}
	r.AuthToken, err = r.Session.GetToken()
	return
}

// Authenticate compares an input password with the hashed password stored in the User model
func (r *Auth) Authenticate(user *User, checkPassword string) (err error) {
	if err = r.authenticate(user, checkPassword); err != nil {
		return
	}
	r.User = user
	if err = r.NewSession(); err != nil {
		return
	}
	err = r.NewToken()
	return
}
