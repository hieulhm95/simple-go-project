package mongodb

import (
	"context"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	CollectionProfile = "collectionProfile"
	CollectionUser    = "collectionUser"
	CollectionPost    = "collectionPost"
)

type Storage struct {
	db      *mongo.Database
	timeout time.Duration
}

func MustStorage(uri, dbName string) *Storage {
	storage, err := NewStorage(uri, dbName)
	if err != nil {
		panic("setup mongo:" + err.Error())
	}

	log.Info("init mongo connection")
	return storage
}

func NewStorage(uri, dbName string) (*Storage, error) {
	client, err := NewClient(uri)
	if err != nil {
		return nil, err
	}

	if err := CreateIndexesCollectionUser(context.TODO(), client.Database(dbName)); err != nil {
		return nil, errors.Wrap(err, "CreateIndexesCollectionUser")
	}

	if err := CreateIndexesCollectionPost(context.TODO(), client.Database(dbName)); err != nil {
		return nil, errors.Wrap(err, "CreateIndexesCollectionPost")
	}

	if err := CreateIndexesCollectionProfile(context.TODO(), client.Database(dbName)); err != nil {
		return nil, errors.Wrap(err, "CreateIndexesCollectionProfile")
	}

	return &Storage{db: client.Database(dbName), timeout: 3 * time.Second}, nil
}

func NewClient(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	clientOptions = clientOptions.SetDirect(true).SetReadPreference(readpref.Primary())

	ctx := context.TODO()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, errors.Wrap(err, "mongo.Connect")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "client.Ping")
	}

	return client, nil
}
