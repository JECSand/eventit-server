package auth

import (
	"errors"
	"github.com/JECSand/eventit-server/domains/shared/enums"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"time"
)

// Session stores the structured data from a session token for use
type Session struct {
	ProfileId string     `json:"profileId,omitempty"`
	Role      enums.Role `json:"role,omitempty"`
}

func NewSession(profileId string, role enums.Role) *Session {
	return &Session{
		profileId,
		role,
	}
}

func LoadSession(tokenString string) (*Session, error) {
	if tokenString == "" {
		return &Session{}, errors.New("no token provided")
	}
	c, e := DecodeJWT(tokenString)
	if e != nil {
		return &Session{}, e
	}
	return &Session{
		ProfileId: c.ProfileId,
		Role:      c.Role,
	}, nil
}

func (s *Session) GetToken() (string, error) {
	if s.ProfileId == "" || s.Role == 0 {
		c := 0
		errMsg := "cannot GetToken() from a session with the following errors: "
		if s.ProfileId == "" {
			errMsg += "missing profileId"
			c++
		}
		if s.Role == 0 {
			if c > 0 {
				errMsg += " and "
			}
			errMsg += "missing or invalid role"
		}
		return "", errors.New(errMsg)
	}
	secret := viper.GetString("auth_jwt_secret")
	claims := AppClaims{
		s.ProfileId,
		s.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			//Issuer:    "test",
			//Subject:   "somebody",
			//ID:        "1",
			//Audience:  []string{"somebody_else"},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(secret))
	// fmt.Println(tokenString, err)
}
