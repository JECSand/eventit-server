package services

import (
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
	//userRec, err := repos.NewUserRecord(user)
	//if err != nil {
	//	return user, err
	//}
	/*
		userRec, err = us.userRepo.Handler.UpdateOne(&repos.UserRecord{Id: userRec.Id}, userRec)
		if err != nil {
			return user, err
		}
	*/
	return user, nil
}
