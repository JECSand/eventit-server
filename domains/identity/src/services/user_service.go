package services

import (
	"github.com/JECSand/eventit-server/domains/identity/src/models"
	"github.com/JECSand/eventit-server/domains/identity/src/repositories"
)

// UserService is used by the app to manage all user related controllers and functionality
type UserService struct {
	userRepo *repositories.UserRepo
}

// NewUserService is an exported function used to initialize a new UserService struct
func NewUserService(uHandler *repositories.UserRepo) *UserService {
	return &UserService{uHandler}
}

func (us *UserService) Create(user *models.User) (*models.User, error) {
	// TODO ADD LOGIC HERE
	return user, nil
}

func (us *UserService) Update(user *models.User) (*models.User, error) {
	// TODO ADD LOGIC HERE
	return user, nil
}

func (us *UserService) DeleteById(id string) error {
	// TODO ADD LOGIC HERE
	return nil
}

func (us *UserService) FindById(id string) (*models.User, error) {
	var user models.User
	// TODO ADD LOGIC HERE
	return &user, nil
}

func (us *UserService) FindByEmail(email string) (*models.User, error) {
	var user models.User
	// TODO ADD LOGIC HERE
	return &user, nil
}

func (us *UserService) FindMany() (users []*models.User, err error) {
	// TODO ADD LOGIC HERE
	return
}
