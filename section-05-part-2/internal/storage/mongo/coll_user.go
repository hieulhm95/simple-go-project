package mongostore

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"section-05-part-2/internal/entity"
	"section-05-part-2/pkg/hashpass"
	"time"
)

func NewUserCollection(db *mongo.Database, collName string) *userCollection {
	return &userCollection{
		client:  db.Collection(collName),
		timeout: 3 * time.Second,
	}
}

type userCollection struct {
	client  *mongo.Collection
	timeout time.Duration
}

func (c *userCollection) Create(info entity.UserInfo) error {
	doc := NewUserDocument(info)
	ctx, cancelFn := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFn()

	if _, err := c.client.InsertOne(ctx, doc); err != nil {
		return err
	}

	//Create indexes whenever create user successfully
	indexErr := CreateUserCollectionIndex(c)
	if indexErr != nil {
		return indexErr
	}
	return nil
}

func (c *userCollection) ChangePassword(username string, info entity.ChangePasswordRequest) error {
	// oldPass and newPassword should not be duplicate
	ctx, cancelFn := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFn()

	res, err := c.Query(username)

	if err != nil {
		return err
	}
	//match := checkPasswordMatching(info.OldPassword, res.HashPass)
	//
	//if !match {
	//	return fmt.Errorf("Old password is not matching")
	//}

	newPassword := info.NewPassword
	ok := IsPasswordDuplicate(newPassword, res.HashPass)
	if ok {
		return fmt.Errorf("Password is duplicated")
	}
	hassNewPassword := hashpass.HashPassword(newPassword)

	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"hasspass": hassNewPassword}}
	if _, err := c.client.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func checkPasswordMatching(oldPassword, oldHassPassword string) bool {
	return oldHassPassword == hashpass.HashPassword(oldPassword)
}

func IsPasswordDuplicate(newPassword, hashPassword string) bool {
	return hashPassword == hashpass.HashPasswordLogin(newPassword, hashPassword)
}

func (c *userCollection) Query(username string) (entity.UserInfo, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFn()
	result := entity.UserInfo{}
	filter := bson.D{{"username", username}}
	err := c.client.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return entity.UserInfo{}, err
	}

	return result, nil
}
