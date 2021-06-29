package mongodb

import (
	"context"
	log "github.com/sirupsen/logrus"
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

// NewClient returns new Client instance
func NewClient(opt *ClientOptions) *Client {
	return &Client{
		opt: opt,
	}
}

// Client is a bubulearn mongo client
type Client struct {
	opt    *ClientOptions
	client *mongo.Client
	db     *mongo.Database
}

// Init initialized mongo DB connection
func (c *Client) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var err error

	opt := &options.ClientOptions{
		Hosts: c.opt.Hosts,
	}

	if c.opt.User != "" {
		opt.SetAuth(options.Credential{
			AuthSource: c.opt.Database,
			Username:   c.opt.User,
			Password:   c.opt.Password,
		})
	}

	c.client, err = mongo.Connect(ctx, opt)
	if err != nil {
		return err
	}

	err = c.client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	c.db = c.client.Database(c.opt.Database)

	return nil
}

// GetCollection returns mongo.Collection instance by its name
func (c *Client) GetCollection(name string) *mongo.Collection {
	return c.db.Collection(name)
}

// GetDAO returns DAO instance with current Client
func (c *Client) GetDAO() *DAO {
	return &DAO{
		Client: c,
	}
}

// Close closes DB connection context
func (c *Client) Close() {
	if c.client != nil {
		err := c.client.Disconnect(context.Background())
		if err != nil {
			log.Error(err)
		}
	}
}