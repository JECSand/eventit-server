package databases

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBRecord is an abstraction of the db model types
type DBRecord interface {
	toDoc() (doc bson.D, err error)
	bsonFilter() (doc bson.D, err error)
	bsonUpdate() (doc bson.D, err error)
	bsonLoad(doc bson.D) (err error)
	addTimeStamps(newRecord bool)
	addObjectID()
	postProcess() (err error)
	getID() (id interface{})
	update(doc interface{}) (err error)
	match(doc interface{}) bool
}

// DBCollection is an abstraction of the dbClient and testDBClient types
type DBCollection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (cur *mongo.Cursor, err error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

// DBClient is an abstraction of the dbClient and testDBClient types
type DBClient interface {
	Connect() error
	Close() error
	GetBucket(bucketName string) (*gridfs.Bucket, error)
	GetCollection(collectionName string) DBCollection
	NewDBHandler(collectionName string) *DBRepo[DBRecord]
	// NewUserHandler() *DBHandler[*userModel]
	// NewGroupHandler() *DBHandler[*groupModel]
	// NewBlacklistHandler() *DBHandler[*blacklistModel]
	// NewTaskHandler() *DBHandler[*taskModel]
	// NewFileHandler() *DBHandler[*fileModel]
}

// DBCursor is an abstraction of the dbClient and testDBClient types
type DBCursor interface {
	Next(ctx context.Context) bool
	Decode(val interface{}) error
	Close(ctx context.Context) error
}

// checkCursorENV returns a DBCursor based on the ENV
func checkCursorENV(cur *mongo.Cursor) DBCursor {
	// if os.Getenv("ENV") == "test" {
	//  	return newTestMongoCursor(cur)
	// }
	return cur
}
