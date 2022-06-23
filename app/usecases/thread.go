package usecases

import (
	"db_forum/app/models"
	"db_forum/app/repositories"
	"db_forum/pkg"
	"strconv"
)

type ThreadUsecase interface {
	CreatePosts(slugOrID string, posts *models.Posts) (err error)
	Get(slugOrID string) (thread *models.Thread, err error)
	Update(slugOrID string, thread *models.Thread) (err error)
	GetPosts(slugOrID string, limit, since int, sort string, desc bool) (posts *models.Posts, err error)
	Vote(slugOrID string, vote *models.Vote) (thread *models.Thread, err error)
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

func (threadUsecase *ThreadUsecaseImpl) CreatePosts(slugOrID string, posts *models.Posts) (err error) {
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
		if parentPost.Thread != thread.ID {
			err = pkg.ErrParentPostFromOtherThread
			return
		}
	}
	_, err = threadUsecase.repoUser.GetByNickname((*posts)[0].Author)
	if err != nil {
		err = pkg.ErrUserNotFound
		return
	}

	err = threadUsecase.repoThread.CreatePosts(thread, posts)
	return
}

func (threadUsecase *ThreadUsecaseImpl) Get(slugOrID string) (thread *models.Thread, err error) {
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

func (threadUsecase *ThreadUsecaseImpl) Update(slugOrID string, thread *models.Thread) (err error) {
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

	err = threadUsecase.repoThread.Update(oldThread)
	if err != nil {
		return
	}

	*thread = *oldThread

	return
}

func (threadUsecase *ThreadUsecaseImpl) GetPosts(slugOrID string, limit, since int, sort string, desc bool) (posts *models.Posts, err error) {
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
		postsSlice, err = threadUsecase.repoThread.GetPostsTree(thread.ID, limit, since, desc)
	case "parent_tree":
		postsSlice, err = threadUsecase.repoThread.GetPostsParentTree(thread.ID, limit, since, desc)
	default:
		postsSlice, err = threadUsecase.repoThread.GetPostsFlat(thread.ID, limit, since, desc)
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

func (threadUsecase *ThreadUsecaseImpl) Vote(slugOrID string, vote *models.Vote) (thread *models.Thread, err error) {
	id, errConv := strconv.Atoi(slugOrID)

	if errConv != nil {
		thread, err = threadUsecase.repoThread.GetBySlug(slugOrID)
	} else {
		thread, err = threadUsecase.repoThread.GetById(int64(id))
	}

	err = threadUsecase.repoVote.Vote(thread.ID, vote)
	if err != nil {
		err = pkg.ErrUserNotFound
		return
	}
	thread.Votes, err = threadUsecase.repoThread.GetVotes(thread.ID)
	return
}
