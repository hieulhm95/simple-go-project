package mongostore

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"section-05-part-1/internal/entity"
	"time"
)

func NewUserCollection(uri string, dbName, collName string) *userCollection {
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	return &userCollection{
		client:  client.Database(dbName).Collection(collName),
		timeout: 3 * time.Second,
	}
}

type userCollection struct {
	client  *mongo.Collection
	timeout time.Duration
}

func (c *userCollection) Save(info entity.UserInfo) error {
	ctx, cancelFn := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFn()
	_, err := c.client.InsertOne(ctx, info)
	if err != nil {
		return err
	}

	return nil
}

func (c *userCollection) Get(username string) (entity.UserInfo, error) {
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
