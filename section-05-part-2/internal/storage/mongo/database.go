package mongostore

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func MustDatabase(uri string, dbName string) *mongo.Database {
	clientOptions := options.Client().ApplyURI(uri)
	// #TODO 1: Config read preference option from primary
	primary := readpref.Primary()
	clientOptions = clientOptions.SetReadPreference(primary)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	return client.Database(dbName)
}
