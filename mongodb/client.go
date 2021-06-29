package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

// ClientOptions is a bubulearn mongo client options
type ClientOptions struct {
	Hosts    []string
	Database string
	User     string
	Password string
}

// CreateMongoClient creates new mongo client and database instances
func CreateMongoClient(opt *ClientOptions) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var err error

	mOpt := &options.ClientOptions{
		Hosts: opt.Hosts,
	}

	if opt.User != "" {
		mOpt.SetAuth(options.Credential{
			AuthSource: opt.Database,
			Username:   opt.User,
			Password:   opt.Password,
		})
	}

	client, err := mongo.Connect(ctx, mOpt)
	if err != nil {
		return nil, nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, nil, err
	}

	db := client.Database(opt.Database)

	return client, db, nil
}
