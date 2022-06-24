package usecases

import (
	"db_forum/app/models"
	"db_forum/app/repositories"
	"db_forum/pkg"
	"strconv"
)

type ThreadUsecase interface {
	CreateNewPosts(slugOrID string, posts *models.Posts) (err error)
	GetInfoAboutThread(slugOrID string) (thread *models.Thread, err error)
	UpdateThread(slugOrID string, thread *models.Thread) (err error)
	GetThreadPosts(slugOrID string, limit, since int, sort string, desc bool) (posts *models.Posts, err error)
	VoteForThread(slugOrID string, vote *models.Vote) (thread *models.Thread, err error)
}

type ThreadUsecaseImpl struct {
	repoVote   repositories.VoteRepository
	repoThread repositories.ThreadRepository
	repoUser   repositories.UserRepository
	repoPost   repositories.PostRepository
}

func MakeThreadUseCase(vote repositories.VoteRepository, thread repositories.ThreadRepository,
	user repositories.UserRepository, post repositories.PostRepository) ThreadUsecase {
	return &ThreadUsecaseImpl{repoVote: vote, repoThread: thread, repoUser: user, repoPost: post}
}

func (threadUsecase *ThreadUsecaseImpl) CreateNewPosts(slugOrID string, posts *models.Posts) (err error) {
	var thread *models.Thread
	id, errConv := strconv.Atoi(slugOrID)
	if errConv != nil {
		thread, err = threadUsecase.repoThread.GetBySlug(slugOrID)
	} else {
		thread, err = threadUsecase.repoThread.GetById(int64(id))
	}

	if err != nil {
		err = pkg.ErrThreadNotFound
		return
	}

	if len(*posts) == 0 {
		return
	}

	if (*posts)[0].Parent != 0 {
		var parentPost *models.Post
		parentPost, err = threadUsecase.repoPost.GetPost((*posts)[0].Parent)
		if parentPost.Thread != thread.Id {
			err = pkg.ErrParentPostFromOtherThread
			return
		}
	}
	_, err = threadUsecase.repoUser.GetInfoAboutUser((*posts)[0].Author)
	if err != nil {
		err = pkg.ErrUserNotFound
		return
	}

	err = threadUsecase.repoThread.CreateThreadPosts(thread, posts)
	return
}

func (threadUsecase *ThreadUsecaseImpl) GetInfoAboutThread(slugOrID string) (thread *models.Thread, err error) {
	id, errConv := strconv.Atoi(slugOrID)
	if errConv != nil {
		thread, err = threadUsecase.repoThread.GetBySlug(slugOrID)
	} else {
		thread, err = threadUsecase.repoThread.GetById(int64(id))
	}
	if err != nil {
		err = pkg.ErrThreadNotFound
		return
	}
	return
}

func (threadUsecase *ThreadUsecaseImpl) UpdateThread(slugOrID string, thread *models.Thread) (err error) {
	id, errConv := strconv.Atoi(slugOrID)
	var oldThread *models.Thread
	if errConv != nil {
		oldThread, err = threadUsecase.repoThread.GetBySlug(slugOrID)
	} else {
		oldThread, err = threadUsecase.repoThread.GetById(int64(id))
	}

	if err != nil {
		err = pkg.ErrThreadNotFound
		return
	}

	if thread.Title != "" {
		oldThread.Title = thread.Title
	}
	if thread.Message != "" {
		oldThread.Message = thread.Message
	}

	err = threadUsecase.repoThread.UpdateThread(oldThread)
	if err != nil {
		return
	}

	*thread = *oldThread

	return
}

func (threadUsecase *ThreadUsecaseImpl) GetThreadPosts(slugOrID string, limit, since int, sort string, desc bool) (posts *models.Posts, err error) {
	id, errConv := strconv.Atoi(slugOrID)
	var thread *models.Thread
	if errConv != nil {
		thread, err = threadUsecase.repoThread.GetBySlug(slugOrID)
	} else {
		thread, err = threadUsecase.repoThread.GetById(int64(id))
	}

	if err != nil {
		err = pkg.ErrThreadNotFound
		return
	}

	postsSlice := new([]models.Post)
	switch sort {
	case "tree":
		postsSlice, err = threadUsecase.repoThread.GetThreadPostsTree(thread.Id, limit, since, desc)
	case "parent_tree":
		postsSlice, err = threadUsecase.repoThread.GetThreadPostsParentTree(thread.Id, limit, since, desc)
	default:
		postsSlice, err = threadUsecase.repoThread.GetThreadPostsFlat(thread.Id, limit, since, desc)
	}
	if err != nil {
		return
	}
	posts = new(models.Posts)
	if len(*postsSlice) == 0 {
		*posts = []models.Post{}
	} else {
		*posts = *postsSlice
	}

	return
}

func (threadUsecase *ThreadUsecaseImpl) VoteForThread(slugOrID string, vote *models.Vote) (thread *models.Thread, err error) {
	id, errConv := strconv.Atoi(slugOrID)

	if errConv != nil {
		thread, err = threadUsecase.repoThread.GetBySlug(slugOrID)
	} else {
		thread, err = threadUsecase.repoThread.GetById(int64(id))
	}

	err = threadUsecase.repoVote.VoteForThread(thread.Id, vote)
	if err != nil {
		err = pkg.ErrUserNotFound
		return
	}
	thread.Votes, err = threadUsecase.repoThread.GetThreadVotes(thread.Id)
	return
}
