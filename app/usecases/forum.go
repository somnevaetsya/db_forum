package usecases

import (
	"db_forum/app/models"
	"db_forum/app/repositories"
	"db_forum/pkg"
)

type ForumUsecase interface {
	CreateForum(forum *models.Forum) (err error)
	GetInfoAboutForum(slug string) (forum *models.Forum, err error)
	CreateForumsThread(thread *models.Thread) (err error)
	GetForumUsers(slug string, limit int, since string, desc bool) (users *models.Users, err error)
	GetForumThreads(slug string, limit int, since string, desc bool) (threads *models.Threads, err error)
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
	user, err := forumUsecase.repoUser.GetInfoAboutUser(forum.User)
	if err != nil {
		err = pkg.ErrUserNotFound
		return
	}

	oldForum, err := forumUsecase.repoForum.GetInfoAboutForum(forum.Slug)
	if oldForum.Slug != "" {
		*forum = *oldForum
		err = pkg.ErrForumAlreadyExists
		return
	}

	forum.User = user.Nickname
	err = forumUsecase.repoForum.CreateForum(forum)
	return
}

func (forumUsecase *ForumUseCaseImpl) GetInfoAboutForum(slug string) (forum *models.Forum, err error) {
	forum, err = forumUsecase.repoForum.GetInfoAboutForum(slug)
	if err != nil {
		err = pkg.ErrForumNotExist
	}
	return
}

func (forumUsecase *ForumUseCaseImpl) CreateForumsThread(thread *models.Thread) (err error) {
	forum, err := forumUsecase.repoForum.GetInfoAboutForum(thread.Forum)
	if err != nil {
		err = pkg.ErrForumOrTheadNotFound
		return
	}

	_, err = forumUsecase.repoUser.GetInfoAboutUser(thread.Author)
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
	err = forumUsecase.repoThread.CreateThread(thread)
	return
}

func (forumUsecase *ForumUseCaseImpl) GetForumUsers(slug string, limit int, since string, desc bool) (users *models.Users, err error) {
	_, err = forumUsecase.repoForum.GetInfoAboutForum(slug)
	if err != nil {
		err = pkg.ErrForumNotExist
		return
	}

	usersSlice, err := forumUsecase.repoForum.GetForumUsers(slug, limit, since, desc)
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

func (forumUsecase *ForumUseCaseImpl) GetForumThreads(slug string, limit int, since string, desc bool) (threads *models.Threads, err error) {
	forum, err := forumUsecase.repoForum.GetInfoAboutForum(slug)
	if err != nil {
		err = pkg.ErrForumNotExist
		return
	}

	threadsSlice, err := forumUsecase.repoForum.GetForumThreads(forum.Slug, limit, since, desc)
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
