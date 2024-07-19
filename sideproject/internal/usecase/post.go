package usecase

import (
	"context"
	"io"
	"sideproject/internal/entity"
	"sideproject/internal/repository/kafka"
)

type PostUC interface {
	InsertPost(ctx context.Context, req *entity.CreatePostRequest) (*entity.Post, error)
	UploadPostImage(ctx context.Context, req *entity.UploadImageRequest) (*entity.UploadImageResponse, error)
	GetPostById(ctx context.Context, id string) (*entity.Post, error)
	GetPostList(ctx context.Context, userId string, offset int64, limit int64) (*entity.GetPostListResponse, error)
	LikePost(ctx context.Context, req *entity.LikePostRequest, post *entity.Post) (*entity.Post, error)
}

type ImageBucket interface {
	SaveImage(ctx context.Context, name string, r io.Reader) (string, error)
}

type newPostUC struct {
	repo      PostUC
	gcs       ImageBucket
	pub       kafka.Publisher
	profileUC ProfileUC
}

func NewPostUC(repo PostUC, pub kafka.Publisher, gcs ImageBucket, uc ProfileUC) *newPostUC {
	return &newPostUC{
		repo:      repo,
		gcs:       gcs,
		pub:       pub,
		profileUC: uc,
	}
}

func (uc *newPostUC) GetPostById(ctx context.Context, id string) (*entity.Post, error) {
	resp, err := uc.repo.GetPostById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entity.Post{
		Id:      resp.Id,
		Caption: resp.Caption,
		Image:   resp.Image,
		Likes:   resp.Likes,
	}, nil
}

func (uc *newPostUC) GetPostList(ctx context.Context, userId string, offset int64, limit int64) (*entity.GetPostListResponse, error) {
	resp, err := uc.repo.GetPostList(ctx, userId, offset, limit)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (uc *newPostUC) InsertPost(ctx context.Context, req *entity.CreatePostRequest) (*entity.Post, error) {
	resp, err := uc.repo.InsertPost(ctx, req)
	if err != nil {
		return nil, err
	}
	return &entity.Post{
		Id:      resp.Id,
		Caption: resp.Caption,
		Image:   resp.Image,
		Likes:   nil,
	}, nil
}

func (uc *newPostUC) UploadPostImage(ctx context.Context, req *entity.UploadImageRequest) (*entity.UploadImageResponse, error) {
	imageName, err := uc.gcs.SaveImage(context.TODO(), req.FileName, req.Content)
	if err != nil {
		return nil, err
	}
	return &entity.UploadImageResponse{ImageUrl: imageName}, nil
}

func (uc *newPostUC) LikePost(ctx context.Context, req *entity.LikePostRequest, post *entity.Post) (*entity.Post, error) {
	var foundUserLike string
	var returnData *entity.Post
	userId := req.UserId
	_, err := uc.profileUC.GetProfileByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	for _, item := range post.Likes {
		if item == userId {
			foundUserLike = item
		}
	}
	if foundUserLike != "" {
		var result []string
		for _, value := range post.Likes {
			if value != foundUserLike {
				result = append(result, value)
			}
		}
		post.Likes = result
		resp, err := uc.repo.LikePost(ctx, req, post)
		if err != nil {
			return nil, err
		}
		returnData = resp
	} else {
		post.Likes = append(post.Likes, userId)
		resp, err := uc.repo.LikePost(ctx, req, post)
		if err != nil {
			return nil, err
		}
		if err := uc.pub.PublishLikeNotification(ctx, resp.Likes); err != nil {
			return nil, err
		}
		returnData = resp
	}
	return returnData, nil
}
