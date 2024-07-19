package mongostore

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"section-05-part-1/internal/entity"
	"time"
)

func NewImageCollection(uri string, dbName, collName string) *imageCollection {
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	return &imageCollection{
		client:  client.Database(dbName).Collection(collName),
		timeout: 3 * time.Second,
	}
}

type imageCollection struct {
	client  *mongo.Collection
	timeout time.Duration
}

func (c *imageCollection) Save(info entity.ImageInfo) error {
	doc := NewImageDocument(info)

	ctx, cancelFn := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFn()

	_, err := c.client.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	return nil
}

type ImageDoc struct {
	User      string    `json:"user" bson:"user"`
	Name      string    `json:"name" bson:"name"`
	Path      string    `json:"path" bson:"path"`
	CreatedDt time.Time `json:"createdDt" bson:"createdDt"`
}

func NewImageDocument(info entity.ImageInfo) ImageDoc {
	return ImageDoc{
		User:      info.UserName,
		Name:      info.FileName,
		Path:      info.ImagePath,
		CreatedDt: time.Now(),
	}
}
