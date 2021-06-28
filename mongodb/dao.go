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

// DAO is a mongo collection abstraction
type DAO struct {
	C *mongo.Collection
}

// FetchByID fetches row by ID to the target
func (d *DAO) FetchByID(id string, target interface{}) error {
	ctx, cancel := d.Ctx(1)
	defer cancel()

	filter := bson.M{"_id": id}
	err := d.C.FindOne(ctx, filter).Decode(target)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Error(err)
			return err
		}
		return bubucore.ErrNotFound
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

	cur, err := d.C.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Error(err)
		}
	}(cur, ctx)

	err = cur.All(ctx, &target)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return bubucore.ErrNotFound
		}
		return err
	}

	return nil
}

// Ctx creates new timeout context
func (d *DAO) Ctx(seconds uint) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
}
