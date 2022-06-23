package usecases

import (
	"db_forum/app/models"
	"db_forum/app/repositories"
	"db_forum/pkg"
)

type ForumUsecase interface {
	CreateForum(forum *models.Forum) (err error)
	Get(slug string) (forum *models.Forum, err error)
	CreateThread(thread *models.Thread) (err error)
	GetUsers(slug string, limit int, since string, desc bool) (users *models.Users, err error)
	GetThreads(slug string, limit int, since string, desc bool) (threads *models.Threads, err error)
}

type ForumUseCaseImpl struct {
	repoForum  repositories.ForumRepository
	repoThread repositories.ThreadRepository
	repoUser   repositories.UserRepository
}

func MakeForumUseCase(forum repositories.ForumRepository, thread repositories.ThreadRepository, user repositories.UserRepository) *ForumUseCaseImpl {
	return &ForumUseCaseImpl{repoForum: forum, repoThread: thread, repoUser: user}
}

func (forumUsecase *ForumUseCaseImpl) CreateForum(forum *models.Forum) (err error) {
	user, err := forumUsecase.repoUser.GetByNickname(forum.User)
	if err != nil {
		err = pkg.ErrUserNotFound
		return
	}

	oldForum, err := forumUsecase.repoForum.GetBySlug(forum.Slug)
	if oldForum.Slug != "" {
		*forum = *oldForum
		err = pkg.ErrForumAlreadyExists
		return
	}

	forum.User = user.Nickname
	err = forumUsecase.repoForum.Create(forum)
	return
}

func (forumUsecase *ForumUseCaseImpl) Get(slug string) (forum *models.Forum, err error) {
	forum, err = forumUsecase.repoForum.GetBySlug(slug)
	if err != nil {
		err = pkg.ErrForumNotExist
	}
	return
}

func (forumUsecase *ForumUseCaseImpl) CreateThread(thread *models.Thread) (err error) {
	forum, err := forumUsecase.repoForum.GetBySlug(thread.Forum)
	if err != nil {
		err = pkg.ErrForumOrTheadNotFound
		return
	}

	_, err = forumUsecase.repoUser.GetByNickname(thread.Author)
	if err != nil {
		err = pkg.ErrForumOrTheadNotFound
		return
	}

	oldThread, err := forumUsecase.repoThread.GetBySlug(thread.Slug)
	if oldThread.Slug != "" {
		*thread = *oldThread
		err = pkg.ErrThreadAlreadyExists
		return
	}

	thread.Forum = forum.Slug
	err = forumUsecase.repoThread.Create(thread)
	return
}

func (forumUsecase *ForumUseCaseImpl) GetUsers(slug string, limit int, since string, desc bool) (users *models.Users, err error) {
	_, err = forumUsecase.repoForum.GetBySlug(slug)
	if err != nil {
		err = pkg.ErrForumNotExist
		return
	}

	usersSlice, err := forumUsecase.repoForum.GetUsers(slug, limit, since, desc)
	if err != nil {
		return
	}
	users = new(models.Users)
	if len(*usersSlice) == 0 {
		*users = []models.User{}
	} else {
		*users = *usersSlice
	}

	return
}

func (forumUsecase *ForumUseCaseImpl) GetThreads(slug string, limit int, since string, desc bool) (threads *models.Threads, err error) {
	forum, err := forumUsecase.repoForum.GetBySlug(slug)
	if err != nil {
		err = pkg.ErrForumNotExist
		return
	}

	threadsSlice, err := forumUsecase.repoForum.GetThreads(forum.Slug, limit, since, desc)
	if err != nil {
		return
	}
	threads = new(models.Threads)
	if len(*threadsSlice) == 0 {
		*threads = []models.Thread{}
	} else {
		*threads = *threadsSlice
	}

	return
}
