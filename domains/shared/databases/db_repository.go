package databases

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

// DBRepo is a Generic type struct for organizing dbModel methods
type DBRepo[T DBRecord] struct {
	db         DBClient
	collection DBCollection
}

// FindOne is used to get a dbModel from the db with custom filter
func (h *DBRepo[T]) FindOne(filter T) (T, error) {
	var m T
	f, err := filter.bsonFilter()
	if err != nil {
		return filter, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = h.collection.FindOne(ctx, f).Decode(&m)
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
	f, err := filter.bsonFilter()
	if err != nil {
		return m, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var cur *mongo.Cursor
	if len(f) > 0 {
		cur, err = h.collection.Find(ctx, f)
	} else {
		cur, err = h.collection.Find(ctx, bson.M{})
	}
	if err != nil {
		return m, err
	}
	cursor := checkCursorENV(cur)
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var md T
		cursor.Decode(&md)
		err = md.postProcess()
		if err != nil {
			return m, err
		}
		m = append(m, md)
	}
	return m, nil
}

// UpdateOne Function to update a dbModel from datasource with custom filter and update model
func (h *DBRepo[T]) UpdateOne(filter T, m T) (T, error) {
	f, err := filter.bsonFilter()
	if err != nil {
		return m, err
	}
	m.addTimeStamps(false)
	update, err := m.bsonUpdate()
	if err != nil {
		return m, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err = h.collection.UpdateOne(ctx, f, update)
	if err != nil {
		return m, err
	}
	err = m.postProcess()
	return m, err
}

// InsertOne adds a new dbModel record to a collection
func (h *DBRepo[T]) InsertOne(m T) (T, error) {
	m.addTimeStamps(true)
	m.addObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := h.collection.InsertOne(ctx, m)
	if err != nil {
		return m, err
	}
	err = m.postProcess()
	return m, err
}

// DeleteOne adds a new dbModel record to a collection
func (h *DBRepo[T]) DeleteOne(filter T) (T, error) { //TODO: to be replaced with "soft delete"
	var m T
	f, err := filter.bsonFilter()
	if err != nil {
		return m, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = h.collection.FindOneAndDelete(ctx, f).Decode(&m)
	return m, err
}

// DeleteMany adds a new dbModel record to a collection
func (h *DBRepo[T]) DeleteMany(filter T) (T, error) { //TODO: to be replaced with "soft delete"
	var m T
	f, err := filter.bsonFilter()
	if err != nil {
		return m, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = h.collection.DeleteMany(ctx, f)
	return filter, err
}

// newRoutine returns a new Routine for executing ASYNC DB statements
func (h *DBRepo[T]) newRoutine() *dbRoutine[T] {
	return &dbRoutine[T]{handler: h}
}
