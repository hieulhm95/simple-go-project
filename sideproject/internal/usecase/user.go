package usecase

import (
	"context"
	"fmt"
	"sideproject/internal/entity"
	"sideproject/internal/repository/kafka"
)

type UserUC interface {
	CreateUser(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error)
	GetUser(ctx context.Context, req *entity.UserInfo) (*entity.UserInfo, error)
}

type newUserUC struct {
	repo UserUC
}

func NewUserUC(repo UserUC, pub kafka.Publisher) *newUserUC {
	return &newUserUC{
		repo: repo,
	}
}

func (uc *newUserUC) CreateUser(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error) {
	user, err := uc.repo.CreateUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Create user: %w", err)

	}
	return user, nil
}

func (uc *newUserUC) GetUser(ctx context.Context, req *entity.UserInfo) (*entity.UserInfo, error) {
	user, err := uc.repo.GetUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Login: %w", err)
	}

	return user, nil
}
