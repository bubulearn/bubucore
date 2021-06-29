package mongodb

import (
	"context"
	"github.com/bubulearn/bubucore"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// DAOInterface is an interface for DAOs
type DAOInterface interface {
	FetchByID(id string, target interface{}) error
	FetchAll(target interface{}, opts ...*options.FindOptions) error
	FetchAllF(target interface{}, filter interface{}, opts ...*options.FindOptions) error
}

// NewDAO creates new DAO instance with the specified collection
func NewDAO(collection *mongo.Collection) *DAO {
	return &DAO{
		c: collection,
	}
}

// DAO is a mongo collection abstraction
type DAO struct {
	c *mongo.Collection
}

// FetchByID fetches row by ID to the target
func (d *DAO) FetchByID(id string, target interface{}) error {
	ctx, cancel := d.Ctx(1)
	defer cancel()

	filter := bson.M{"_id": id}
	err := d.C().FindOne(ctx, filter).Decode(target)

	if err != nil {
		return d.Err(err)
	}

	return nil
}

// FetchAll fetches all rows from cursor to the target
func (d *DAO) FetchAll(target interface{}, opts ...*options.FindOptions) error {
	return d.FetchAllF(target, bson.M{}, opts...)
}

// FetchAllF fetches all rows from cursor with filter to the target
func (d *DAO) FetchAllF(target interface{}, filter interface{}, opts ...*options.FindOptions) error {
	ctx, cancel := d.Ctx(10)
	defer cancel()

	cur, err := d.C().Find(ctx, filter, opts...)
	if err != nil {
		return d.Err(err)
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		_ = d.Err(err)
	}(cur, ctx)

	err = cur.All(ctx, &target)
	if err != nil {
		return d.Err(err)
	}

	return nil
}

// C returns Collection
func (d *DAO) C() *mongo.Collection {
	return d.c
}

// Ctx creates new timeout context
func (d *DAO) Ctx(seconds uint) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
}

// Err transforms and log an error if needed
func (d *DAO) Err(err error) error {
	if err == nil {
		return nil
	}
	needLog := true
	switch err {
	case mongo.ErrNoDocuments:
		return bubucore.ErrNotFound
	case bubucore.ErrNotFound:
		needLog = false
	}
	if needLog {
		cName := "_unknown_"
		c := d.C()
		if c != nil {
			cName = c.Name()
		}
		log.WithField("dao", cName).Error(err)
	}
	return err
}
