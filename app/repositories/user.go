package repositories

import (
	"db_forum/app/models"
	"db_forum/pkg/queries"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

type UserRepository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	GetByNickname(nickname string) (*models.User, error)
	GetSimilar(user *models.User) (*[]models.User, error)
}

type UserRepositoryImpl struct {
	db *pgx.ConnPool
}

func MakeUserRepository(db *pgx.ConnPool) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (userRepository *UserRepositoryImpl) Create(user *models.User) error {
	_, err := userRepository.db.Exec(queries.UserCreate, user.Nickname, user.Fullname, user.About, user.Email)
	return err
}

func (userRepository *UserRepositoryImpl) Update(user *models.User) error {
	return userRepository.db.QueryRow(queries.UserUpdate, user.Fullname, user.About, user.Email, user.Nickname).Scan(&user.Fullname, &user.About, &user.Email)
}

func (userRepository *UserRepositoryImpl) GetByNickname(nickname string) (*models.User, error) {
	user := new(models.User)
	err := userRepository.db.QueryRow(queries.UserGet, nickname).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	return user, err
}

func (userRepository *UserRepositoryImpl) GetSimilar(user *models.User) (*[]models.User, error) {
	var usersSlice []models.User

	resultRows, err := userRepository.db.Query(queries.UserGetSimilar, user.Nickname, user.Email)
	if err != nil {
		return nil, err
	}
	defer resultRows.Close()

	for resultRows.Next() {
		user := models.User{}
		err = resultRows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			return nil, err
		}
		usersSlice = append(usersSlice, user)
	}
	return &usersSlice, nil
}
