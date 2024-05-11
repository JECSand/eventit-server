package services

import (
	"errors"
	"github.com/JECSand/eventit-server/domains/identity/src/models"
	repos "github.com/JECSand/eventit-server/domains/identity/src/repositories"
	"time"
)

// AuthService is used by the app to manage all user related controllers and functionality
type AuthService struct {
	userService *UserService
	blacklist   *repos.BlacklistRepo
}

// NewAuthService is an exported function used to initialize a new UserService struct
func NewAuthService(userService *UserService, blHandler *repos.BlacklistRepo) *AuthService {
	return &AuthService{
		userService,
		blHandler,
	}
}

func (us *AuthService) Login(credentials *models.Credentials) (*models.Auth, error) {
	auth := &models.Auth{CreatedAt: time.Now().UTC()}
	if credentials.Password == "" {
		return auth, errors.New("password is empty")
	}
	if credentials.Email == "" {
		return auth, errors.New("email is empty")
	}
	foundUser, err := us.userService.FindByEmail(credentials.Email)
	if err != nil {
		return auth, err
	}
	if err = auth.Authenticate(foundUser, credentials.Password); err != nil {
		return auth, err
	}
	return auth, nil
}

func (us *AuthService) Logout(auth *models.Auth) error {
	if auth.AuthToken == "" {
		return errors.New("token is empty")
	}
	_, err := us.blacklist.Handler.InsertOne(&repos.BlacklistRecord{AuthToken: auth.AuthToken})
	if err == nil {
		return err
	}
	auth.Invalidate()
	return nil
}

func (us *AuthService) Validate(auth *models.Auth) error {
	if auth.AuthToken == "" {
		return errors.New("token is empty")
	}
	if err := auth.LoadSession(); err != nil {
		return err
	}
	foundUser, err := us.userService.FindById(auth.Session.ProfileId)
	if err != nil {
		return err
	}
	auth.User = foundUser
	return nil
}
