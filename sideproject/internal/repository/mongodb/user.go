package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"sideproject/internal/entity"
)

func (s *Storage) GetUser(ctx context.Context, req *entity.UserInfo) (*entity.UserInfo, error) {
	var doc UserDoc
	ctx, cancelFn := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFn()

	filter := bson.M{"username": req.Username}
	if err := s.db.Collection(CollectionUser).FindOne(ctx, filter).Decode(&doc); err != nil {
		return nil, err
	}

	return &entity.UserInfo{
		Id:       doc.Id,
		Username: doc.Username,
		FullName: doc.FullName,
	}, nil
}

func (s *Storage) CreateUser(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error) {
	userDoc := NewUserDocument(req)
	ctx, cancelFn := context.WithTimeout(context.Background(), s.timeout)
	defer cancelFn()

	_, err := s.db.Collection(CollectionUser).InsertOne(ctx, userDoc)
	if err != nil {
		return nil, err
	}
	return &entity.RegisterResponse{
		Username: userDoc.Username,
		FullName: userDoc.FullName,
		UserId:   userDoc.Id,
	}, nil
}
