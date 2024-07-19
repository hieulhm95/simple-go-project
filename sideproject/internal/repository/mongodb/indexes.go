package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexesCollectionUser(ctx context.Context, db *mongo.Database) error {
	c := db.Collection(CollectionUser)

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := c.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	return nil
}

func CreateIndexesCollectionProfile(ctx context.Context, db *mongo.Database) error {
	c := db.Collection(CollectionProfile)

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"id": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := c.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	return nil
}

func CreateIndexesCollectionPost(ctx context.Context, db *mongo.Database) error {
	c := db.Collection(CollectionPost)

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"id": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := c.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	return nil
}
