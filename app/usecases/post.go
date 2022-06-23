package usecases

import (
	"db_forum/app/models"
	"db_forum/app/repositories"
	"db_forum/pkg"
)

type PostUsecase interface {
	GetPost(id int64, relatedData *[]string) (postFull *models.PostFull, err error)
	UpdatePost(post *models.Post) (err error)
}

type PostUsecaseImpl struct {
	repoForum  repositories.ForumRepository
	repoThread repositories.ThreadRepository
	repoUser   repositories.UserRepository
	repoPost   repositories.PostRepository
}

func MakePostUseCase(forum repositories.ForumRepository, thread repositories.ThreadRepository,
	user repositories.UserRepository, post repositories.PostRepository) PostUsecase {
	return &PostUsecaseImpl{repoForum: forum, repoThread: thread, repoUser: user, repoPost: post}
}

func (postUsecase *PostUsecaseImpl) GetPost(id int64, relatedData *[]string) (postFull *models.PostFull, err error) {
	postFull = new(models.PostFull)
	var post *models.Post
	post, err = postUsecase.repoPost.GetPost(id)
	if err != nil {
		err = pkg.ErrPostNotFound
	}
	postFull.Post = post

	for _, data := range *relatedData {
		switch data {
		case "user":
			var author *models.User
			author, err = postUsecase.repoUser.GetByNickname(postFull.Post.Author)
			if err != nil {
				err = pkg.ErrUserNotFound
			}
			postFull.Author = author
		case "forum":
			var forum *models.Forum
			forum, err = postUsecase.repoForum.GetBySlug(postFull.Post.Forum)
			if err != nil {
				err = pkg.ErrForumNotExist
			}
			postFull.Forum = forum
		case "thread":
			var thread *models.Thread
			thread, err = postUsecase.repoThread.GetById(postFull.Post.Thread)
			if err != nil {
				err = pkg.ErrThreadNotFound
			}
			postFull.Thread = thread
		}
	}
	return
}

func (postUsecase *PostUsecaseImpl) UpdatePost(post *models.Post) (err error) {
	oldPost, err := postUsecase.repoPost.GetPost(post.ID)
	if err != nil {
		err = pkg.ErrThreadNotFound
		return
	}

	if post.Message != "" {
		if oldPost.Message != post.Message {
			oldPost.IsEdited = true
		}
		oldPost.Message = post.Message

		err = postUsecase.repoPost.UpdatePost(oldPost)
		if err != nil {
			return
		}
	}

	*post = *oldPost

	return
}
