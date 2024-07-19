package usecase

import (
	"context"
	"io"
	"section-04-part-2/internal/entity"
	"section-04-part-2/pkg/auth"
	"time"
)

type UserStore interface {
	Save(info entity.UserInfo) error
	Get(username string) (entity.UserInfo, error)
}

type ImageBucket interface {
	SaveImage(ctx context.Context, name string, r io.Reader) (string, error)
}

func NewUseCase(userStore UserStore, imageBucket ImageBucket) *ucImplement {
	return &ucImplement{
		userStore: userStore,
		imgBucket: imageBucket,
	}
}

type ucImplement struct {
	userStore UserStore
	imgBucket ImageBucket
}

func (uc *ucImplement) Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error) {
	if err := uc.userStore.Save(entity.UserInfo{
		Username: req.Username,
		Password: req.Password,
		FullName: req.FullName,
		Address:  req.Address,
	}); err != nil {
		return nil, err
	}

	return &entity.RegisterResponse{UserId: req.Username}, nil
}

func (uc *ucImplement) Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error) {
	user, err := uc.userStore.Get(req.UserName)
	if err != nil {
		return nil, err
	}
	if req.Password != user.Password {
		return nil, err
	}
	token, err := auth.GenerateToken(user.Username, 24*time.Second)
	if err != nil {
		return nil, err
	}
	return &entity.LoginResponse{Token: token}, nil
}

func (uc *ucImplement) Self(ctx context.Context, userName string) (*entity.SelfResponse, error) {
	user, err := uc.userStore.Get(userName)
	if err != nil {
		return nil, err
	}
	return &entity.SelfResponse{
		Username: user.Username,
		FullName: user.FullName,
		Address:  user.Address,
	}, nil
}

func (uc *ucImplement) UploadImage(ctx context.Context, req *entity.UploadImageRequest) (*entity.UploadImageResponse, error) {
	imageUrl, err := uc.imgBucket.SaveImage(ctx, req.FileName, req.File)
	if err != nil {
		return nil, err
	}
	return &entity.UploadImageResponse{Url: imageUrl}, nil
}
