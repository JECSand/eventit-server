package services

import (
	"errors"
	"github.com/JECSand/eventit-server/domains/identity/src/models"
	repos "github.com/JECSand/eventit-server/domains/identity/src/repositories"
)

// AuthService is used by the app to manage all user related controllers and functionality
type AuthService struct {
	userService *UserService
	blRepo      *repos.BlacklistRepo
}

// NewAuthService is an exported function used to initialize a new UserService struct
func NewAuthService(userService *UserService, blHandler *repos.BlacklistRepo) *AuthService {
	return &AuthService{
		userService,
		blHandler,
	}
}

func (us *AuthService) Login(user *models.User) (*models.User, error) {
	if user.Password == "" {
		return user, errors.New("password is empty")
	}
	if user.Email == "" {
		return user, errors.New("email is empty")
	}
	foundUser, err := us.userService.FindByEmail(user.Email)
	if err != nil {
		return user, err
	}
	if err = foundUser.Authenticate(user.Password); err != nil {
		return user, err
	}
	// TODO - Generate auth token here and change response struct.
	return foundUser, nil
}
