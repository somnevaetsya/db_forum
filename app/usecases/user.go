package usecases

import (
	"db_forum/app/models"
	"db_forum/app/repositories"
	"db_forum/pkg"
)

type UserUsecase interface {
	CreateNewUser(user *models.User) (users *models.Users, err error)
	GetInfoAboutUser(nickname string) (user *models.User, err error)
	UpdateUser(user *models.User) (err error)
}

type UserUsecaseImpl struct {
	repoUser repositories.UserRepository
}

func MakeUserUseCase(user repositories.UserRepository) UserUsecase {
	return &UserUsecaseImpl{repoUser: user}
}

func (userUsecase *UserUsecaseImpl) CreateNewUser(user *models.User) (users *models.Users, err error) {
	usersSlice, err := userUsecase.repoUser.GetSimilarUsers(user)
	if err != nil {
		err = pkg.ErrUserAlreadyExist
		return
	} else if len(*usersSlice) > 0 {
		users = new(models.Users)
		*users = *usersSlice
		err = pkg.ErrUserAlreadyExist
		return
	}

	err = userUsecase.repoUser.CreateUser(user)
	return
}

func (userUsecase *UserUsecaseImpl) GetInfoAboutUser(nickname string) (user *models.User, err error) {
	user, err = userUsecase.repoUser.GetInfoAboutUser(nickname)
	if err != nil {
		err = pkg.ErrUserNotFound
	}
	return
}

func (userUsecase *UserUsecaseImpl) UpdateUser(user *models.User) (err error) {
	oldUser, err := userUsecase.repoUser.GetInfoAboutUser(user.Nickname)
	if oldUser.Nickname == "" {
		err = pkg.ErrUserNotFound
		return
	}
	if oldUser.Fullname != user.Fullname && user.Fullname == "" {
		user.Fullname = oldUser.Fullname
	}
	if oldUser.About != user.About && user.About == "" {
		user.About = oldUser.About
	}
	if oldUser.Email != user.Email && user.Email == "" {
		user.Email = oldUser.Email
	}
	err = userUsecase.repoUser.UpdateUser(user)
	if err != nil {
		err = pkg.ErrUserDataConflict
	}
	return
}
