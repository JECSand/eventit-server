package services

import (
	"context"
	"errors"
	"github.com/JECSand/eventit-server/domains/identity/src/models"
	repos "github.com/JECSand/eventit-server/domains/identity/src/repositories"
	"github.com/JECSand/eventit-server/domains/shared/utilities"
)

// UserService is used by the app to manage all user related controllers and functionality
type UserService struct {
	userRepo *repos.UserRepo
}

// NewUserService is an exported function used to initialize a new UserService struct
func NewUserService(uHandler *repos.UserRepo) *UserService {
	return &UserService{uHandler}
}

func (us *UserService) Create(user *models.User) (*models.User, error) {
	userRec, err := repos.NewUserRecord(user)
	if err != nil {
		return user, err
	}
	userRec, err = us.userRepo.Handler.InsertOne(userRec)
	if err != nil {
		return user, err
	}
	return userRec.ToRoot(), nil
}

func (us *UserService) Update(user *models.User) (*models.User, error) {
	userRec, err := repos.NewUserRecord(user)
	if err != nil {
		return user, err
	}
	userRec, err = us.userRepo.Handler.UpdateOne(&repos.UserRecord{Id: userRec.Id}, userRec)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (us *UserService) DeleteById(id string) error {
	// TODO ADD LOGIC HERE
	userRec, err := repos.NewUserRecord(&models.User{Id: id})
	if err != nil {
		return err
	}
	_, err = us.userRepo.Handler.DeleteOne(userRec)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) findOne(filter *models.User) (user *models.User, err error) {
	var userRec *repos.UserRecord
	userRec, err = repos.NewUserRecord(filter)
	if err != nil {
		return
	}
	userRec, err = us.userRepo.Handler.FindOne(userRec)
	if err == nil {
		user = userRec.ToRoot()
	}
	return
}

func (us *UserService) FindById(id string) (user *models.User, err error) {
	user.Id = id
	user, err = us.findOne(user)
	return
}

func (us *UserService) FindByEmail(email string) (user *models.User, err error) {
	if ok := utilities.IsValidEmail(email); ok {
		user.Email = email
		user, err = us.findOne(user)
		return
	}
	err = errors.New("invalid email")
	return
}

func (us *UserService) Find(ctx context.Context, user *models.User, pagination *utilities.Pagination) (*models.UsersPage, error) {
	userRec, err := repos.NewUserRecord(user)
	if err != nil {
		return &models.UsersPage{}, err
	}
	count, err := us.userRepo.Handler.Count(userRec)
	if err != nil {
		return &models.UsersPage{}, err
	}
	if count == 0 {
		return &models.UsersPage{
			TotalCount: 0,
			TotalPages: 0,
			Page:       0,
			Size:       0,
			HasMore:    false,
			Users:      make([]*models.User, 0),
		}, nil
	}
	userRecs, err := us.userRepo.Handler.PaginatedFind(ctx, userRec, pagination)
	if err != nil {
		return &models.UsersPage{}, err
	}
	users := repos.LoadUserRecords(userRecs)
	return &models.UsersPage{
		TotalCount: count,
		TotalPages: int64(pagination.GetTotalPages(int(count))),
		Page:       int64(pagination.GetPage()),
		Size:       int64(pagination.GetSize()),
		HasMore:    pagination.GetHasMore(int(count)),
		Users:      users,
	}, nil
}
