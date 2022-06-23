package usecases

import (
	"db_forum/app/models"
	"db_forum/app/repositories"
	"db_forum/pkg"
)

type UserUsecase interface {
	Create(user *models.User) (users *models.Users, err error)
	Get(nickname string) (user *models.User, err error)
	Update(user *models.User) (err error)
}

type UserUsecaseImpl struct {
	repoUser repositories.UserRepository
}

func MakeUserUseCase(user repositories.UserRepository) UserUsecase {
	return &UserUsecaseImpl{repoUser: user}
}

func (userUsecase *UserUsecaseImpl) Create(user *models.User) (users *models.Users, err error) {
	usersSlice, err := userUsecase.repoUser.GetSimilar(user)
	if err != nil {
		err = pkg.ErrUserAlreadyExist
		return
	} else if len(*usersSlice) > 0 {
		users = new(models.Users)
		*users = *usersSlice
		err = pkg.ErrUserAlreadyExist
		return
	}

	err = userUsecase.repoUser.Create(user)
	return
}

func (userUsecase *UserUsecaseImpl) Get(nickname string) (user *models.User, err error) {
	user, err = userUsecase.repoUser.GetByNickname(nickname)
	if err != nil {
		err = pkg.ErrUserNotFound
	}
	return
}

func (userUsecase *UserUsecaseImpl) Update(user *models.User) (err error) {
	oldUser, err := userUsecase.repoUser.GetByNickname(user.Nickname)
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
	err = userUsecase.repoUser.Update(user)
	if err != nil {
		err = pkg.ErrUserDataConflict
	}
	return
}
