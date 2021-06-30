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

// DAO is an interface for DAOs
type DAO interface {
	FetchByID(id string, target interface{}, opts ...*options.FindOneOptions) error
	FetchByIDs(ids []string, target interface{}, opts ...*options.FindOptions) error
	FetchByExIDs(ids []string, target interface{}, opts ...*options.FindOptions) error

	FetchOne(target interface{}, filter interface{}, opts ...*options.FindOneOptions) error
	FetchAll(target interface{}, opts ...*options.FindOptions) error
	FetchAllF(target interface{}, filter interface{}, opts ...*options.FindOptions) error

	Ctx(seconds uint) (context.Context, context.CancelFunc)
	Err(err error) error
}

// NewDAOMg creates new DAOMg instance with the specified collection
func NewDAOMg(collection *mongo.Collection) *DAOMg {
	return &DAOMg{
		c: collection,
	}
}

// DAOMg is a mongo collection abstraction
type DAOMg struct {
	c *mongo.Collection
}

// FetchByID fetches row by ID to the target
func (d *DAOMg) FetchByID(id string, target interface{}, opts ...*options.FindOneOptions) error {
	ctx, cancel := d.Ctx(1)
	defer cancel()

	filter := bson.M{"_id": id}
	err := d.C().FindOne(ctx, filter, opts...).Decode(target)

	if err != nil {
		return d.Err(err)
	}

	return nil
}

// FetchByIDs fetches rows by IDs list
func (d *DAOMg) FetchByIDs(ids []string, target interface{}, opts ...*options.FindOptions) error {
	filter := bson.M{"_id": bson.M{"$in": ids}}
	return d.FetchAllF(target, filter, opts...)
}

// FetchByExIDs fetches rows by exclude IDs list
func (d *DAOMg) FetchByExIDs(ids []string, target interface{}, opts ...*options.FindOptions) error {
	filter := bson.M{"_id": bson.M{"$nin": ids}}
	return d.FetchAllF(target, filter, opts...)
}

// FetchOne fetches one row by the filter
func (d *DAOMg) FetchOne(target interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
	ctx, cancel := d.Ctx(1)
	defer cancel()

	err := d.C().FindOne(ctx, filter, opts...).Decode(target)

	if err != nil {
		return d.Err(err)
	}

	return nil
}

// FetchAll fetches all rows from cursor to the target
func (d *DAOMg) FetchAll(target interface{}, opts ...*options.FindOptions) error {
	return d.FetchAllF(target, bson.M{}, opts...)
}

// FetchAllF fetches all rows from cursor with filter to the target
func (d *DAOMg) FetchAllF(target interface{}, filter interface{}, opts ...*options.FindOptions) error {
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
func (d *DAOMg) C() *mongo.Collection {
	return d.c
}

// Ctx creates new timeout context
func (d *DAOMg) Ctx(seconds uint) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
}

// Err transforms and log an error if needed
func (d *DAOMg) Err(err error) error {
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
