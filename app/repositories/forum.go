package repositories

import (
	"db_forum/app/models"
	"db_forum/pkg/queries"
	"fmt"

	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

type ForumRepository interface {
	Create(forum *models.Forum) (err error)
	GetBySlug(slug string) (forum *models.Forum, err error)
	GetUsers(slug string, limit int, since string, desc bool) (*[]models.User, error)
	GetThreads(slug string, limit int, since string, desc bool) (threads *[]models.Thread, err error)
}

type ForumRepositoryImpl struct {
	db *pgx.ConnPool
}

func MakeForumRepository(db *pgx.ConnPool) ForumRepository {
	return &ForumRepositoryImpl{db: db}
}

func (forumRepository *ForumRepositoryImpl) Create(forum *models.Forum) (err error) {
	_, err = forumRepository.db.Exec(queries.ForumCreate, forum.Title, forum.User, forum.Slug)
	return err
}

func (forumRepository *ForumRepositoryImpl) GetBySlug(slug string) (forum *models.Forum, err error) {
	forum = new(models.Forum)
	err = forumRepository.db.QueryRow(queries.ForumGetBySlug, slug).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	fmt.Println("GET FORUM", forum, err)
	return forum, err
}

func (forumRepository *ForumRepositoryImpl) GetUsers(slug string, limit int, since string, desc bool) (users *[]models.User, err error) {
	var bufUser []models.User

	var query string

	var result *pgx.Rows
	var innerError error

	if since != "" {
		if desc {
			query = queries.ForumGetUsersSinceDesc
		} else {
			query = queries.ForumGetUsersSince
		}
		result, innerError = forumRepository.db.Query(query, slug, since, limit)
		if innerError != nil {
			return
		}
	} else {
		if desc {
			query = queries.ForumGetUsersDesc
		} else {
			query = queries.ForumGetUsers
		}
		result, innerError = forumRepository.db.Query(query, slug, limit)
		if innerError != nil {
			return
		}
	}

	defer result.Close()

	for result.Next() {
		user := models.User{}
		err = result.Scan(
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email)
		if err != nil {
			return
		}
		bufUser = append(bufUser, user)
	}
	return &bufUser, nil
}

func (forumRepository *ForumRepositoryImpl) GetThreads(slug string, limit int, since string, desc bool) (threads *[]models.Thread, err error) {
	var bufThreads []models.Thread

	var query string

	var result *pgx.Rows
	var innerError error

	if since != "" {
		if desc {
			query = queries.ForumGetThreadsSinceDesc
		} else {
			query = queries.ForumGetThreadsSince
		}
		result, innerError = forumRepository.db.Query(query, slug, since, limit)
		if innerError != nil {
			return
		}
	} else {
		if desc {
			query = queries.ForumGetThreadsDesc
		} else {
			query = queries.ForumGetThreads
		}
		result, innerError = forumRepository.db.Query(query, slug, limit)
		if innerError != nil {
			return
		}
	}

	if err != nil {
		return
	}
	defer result.Close()

	for result.Next() {
		thread := models.Thread{}
		err = result.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created)
		if err != nil {
			return
		}
		bufThreads = append(bufThreads, thread)
	}
	return &bufThreads, nil
}
