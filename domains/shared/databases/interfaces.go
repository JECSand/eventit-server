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
	ToDoc() (doc bson.D, err error)
	BsonFilter() (doc bson.D, err error)
	BsonUpdate() (doc bson.D, err error)
	BsonLoad(doc bson.D) (err error)
	AddTimeStamps(newRecord bool)
	AddObjectID()
	PostProcess() (err error)
	GetID() (id interface{})
	Update(doc interface{}) (err error)
	Match(doc interface{}) bool
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
	Err() error
}

// checkCursorENV returns a DBCursor based on the ENV
func checkCursorENV(cur *mongo.Cursor) DBCursor {
	// if os.Getenv("ENV") == "test" {
	//  	return newTestMongoCursor(cur)
	// }
	return cur
}
