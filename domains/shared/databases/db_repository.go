package databases

import (
	"context"
	"github.com/JECSand/eventit-server/domains/shared/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

// DBRepo is a Generic type struct for organizing dbModel methods
type DBRepo[T DBRecord] struct {
	DB         DBClient
	Collection DBCollection
}

// FindOne is used to get a dbModel from the db with custom filter
func (h *DBRepo[T]) FindOne(filter T) (T, error) {
	var m T
	f, err := filter.BsonFilter()
	if err != nil {
		return filter, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = h.Collection.FindOne(ctx, f).Decode(&m)
	if err != nil {
		return filter, err
	}
	return m, nil
}

// FindOneAsync is used to get a dbModel from the db with custom filter
func (h *DBRepo[T]) FindOneAsync(tCh chan T, eCh chan error, filter T, wg *sync.WaitGroup) {
	defer wg.Done()
	t, err := h.FindOne(filter)
	tCh <- t
	eCh <- err
}

// FindMany is used to get a slice of dbModels from the db with custom filter
func (h *DBRepo[T]) FindMany(filter T) ([]T, error) {
	var m []T
	f, err := filter.BsonFilter()
	if err != nil {
		return m, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var cur *mongo.Cursor
	if len(f) > 0 {
		cur, err = h.Collection.Find(ctx, f)
	} else {
		cur, err = h.Collection.Find(ctx, bson.M{})
	}
	if err != nil {
		return m, err
	}
	cursor := checkCursorENV(cur)
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var md T
		cursor.Decode(&md)
		err = md.PostProcess()
		if err != nil {
			return m, err
		}
		m = append(m, md)
	}
	return m, nil
}

// PaginatedFind is used to get a slice of dbModels from the db with custom filter
func (h *DBRepo[T]) PaginatedFind(ctx context.Context, filter T, pagination *utilities.Pagination) ([]T, error) {
	var m []T
	f, err := filter.BsonFilter()
	if err != nil {
		return m, err
	}
	var cur *mongo.Cursor
	limit := int64(pagination.GetLimit())
	skip := int64(pagination.GetOffset())
	if len(f) > 0 {
		cur, err = h.Collection.Find(ctx, f, &options.FindOptions{
			Limit: &limit,
			Skip:  &skip,
		})
	} else {
		cur, err = h.Collection.Find(ctx, bson.M{}, &options.FindOptions{
			Limit: &limit,
			Skip:  &skip,
		})
	}
	if err != nil {
		return m, err
	}
	cursor := checkCursorENV(cur)
	defer cursor.Close(ctx)
	m = make([]T, 0, pagination.GetSize())
	for cursor.Next(ctx) {
		var md T
		if err = cursor.Decode(&md); err != nil {
			return nil, err
		}
		err = md.PostProcess()
		if err != nil {
			return m, err
		}
		m = append(m, md)
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

// Count Function to update a dbModel from datasource with custom filter and update model
func (h *DBRepo[T]) Count(filter T) (int64, error) {
	f, err := filter.BsonFilter()
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return h.Collection.CountDocuments(ctx, f)
}

// UpdateOne Function to update a dbModel from datasource with custom filter and update model
func (h *DBRepo[T]) UpdateOne(filter T, m T) (T, error) {
	f, err := filter.BsonFilter()
	if err != nil {
		return m, err
	}
	m.AddTimeStamps(false)
	update, err := m.BsonUpdate()
	if err != nil {
		return m, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err = h.Collection.UpdateOne(ctx, f, update)
	if err != nil {
		return m, err
	}
	err = m.PostProcess()
	return m, err
}

// InsertOne adds a new dbModel record to a collection
func (h *DBRepo[T]) InsertOne(m T) (T, error) {
	m.AddTimeStamps(true)
	m.AddObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := h.Collection.InsertOne(ctx, m)
	if err != nil {
		return m, err
	}
	err = m.PostProcess()
	return m, err
}

// DeleteOne adds a new dbModel record to a collection
func (h *DBRepo[T]) DeleteOne(filter T) (T, error) { //TODO: to be replaced with "soft delete"
	var m T
	f, err := filter.BsonFilter()
	if err != nil {
		return m, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = h.Collection.FindOneAndDelete(ctx, f).Decode(&m)
	return m, err
}

// DeleteMany adds a new dbModel record to a collection
func (h *DBRepo[T]) DeleteMany(filter T) (T, error) { //TODO: to be replaced with "soft delete"
	var m T
	f, err := filter.BsonFilter()
	if err != nil {
		return m, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = h.Collection.DeleteMany(ctx, f)
	return filter, err
}

// newRoutine returns a new Routine for executing ASYNC DB statements
func (h *DBRepo[T]) newRoutine() *dbRoutine[T] {
	return &dbRoutine[T]{handler: h}
}
