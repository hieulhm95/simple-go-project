package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"sideproject/internal/entity"
)

func (s *Storage) GetProfile(ctx context.Context, id string) (*entity.Profile, error) {
	var doc ProfileDoc
	ctx, cancelFn := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFn()

	filter := bson.M{"id": id}
	if err := s.db.Collection(CollectionProfile).FindOne(ctx, filter).Decode(&doc); err != nil {
		return nil, err
	}

	return &entity.Profile{
		Id:   doc.Id,
		Name: doc.Name,
	}, nil
}

func (s *Storage) GetProfileByUserId(ctx context.Context, userId string) (*entity.Profile, error) {
	var doc ProfileDoc
	ctx, cancelFn := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFn()

	filter := bson.M{"user": userId}
	if err := s.db.Collection(CollectionProfile).FindOne(ctx, filter).Decode(&doc); err != nil {
		return nil, err
	}

	return &entity.Profile{
		Id:   doc.Id,
		Name: doc.Name,
	}, nil
}

func (s *Storage) InsertProfile(ctx context.Context, profile *entity.Profile) (*entity.Profile, error) {
	profileDoc := NewProfileDoc(profile)
	_, err := s.db.Collection(CollectionProfile).InsertOne(ctx, profileDoc)
	if err != nil {
		return nil, err
	}
	return &entity.Profile{
		Id:       profileDoc.Id,
		Name:     profileDoc.Name,
		User:     profileDoc.User,
		Username: profileDoc.Username,
	}, nil
}
