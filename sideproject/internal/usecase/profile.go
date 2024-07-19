package usecase //nolint:revive

import (
	"context"
	"sideproject/internal/entity"
	"sideproject/internal/repository/kafka"
)

type ProfileUC interface {
	GetProfile(ctx context.Context, id string) (*entity.Profile, error)
	InsertProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error)
	GetProfileByUserId(ctx context.Context, userId string) (*entity.Profile, error)
}

type newProfileUC struct {
	repo ProfileUC
}

func NewProfileUC(repo ProfileUC, pub kafka.Publisher) *newProfileUC {
	return &newProfileUC{
		repo: repo,
	}
}

func (uc *newProfileUC) GetProfile(ctx context.Context, id string) (res *entity.Profile, err error) {
	profile, err := uc.repo.GetProfile(ctx, id)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (uc *newProfileUC) GetProfileByUserId(ctx context.Context, userId string) (res *entity.Profile, err error) {
	profile, err := uc.repo.GetProfileByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (uc *newProfileUC) InsertProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error) {
	resp, err := uc.repo.InsertProfile(ctx, profile)
	if err != nil {
		return nil, err

	}
	return resp, nil
}
