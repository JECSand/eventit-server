package databases

import (
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

// dbClient manages a database connection
type dbClient struct {
	connectionURI string
	client        *mongo.Client
}

/*
// InitializeNewClient returns an initialized DBClient based on the ENV
func InitializeNewClient() (DBClient, error) {
	if os.Getenv("ENV") == "test" {
		return initializeNewTestClient()
	}
	return initializeNewClient()
}
*/

// InitializeNewClient is a function that takes a mongoUri string and outputs a connected mongo client for the app to use
func initializeNewClient() (*dbClient, error) {
	newDBClient := dbClient{connectionURI: os.Getenv("MONGO_URI")}
	var err error
	newDBClient.client, err = mongo.NewClient(options.Client().ApplyURI(newDBClient.connectionURI))
	return &newDBClient, err
}

// Connect opens a new connection to the database
func (db *dbClient) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return db.client.Connect(ctx)
}

// Close closes an open DB connection
func (db *dbClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return db.client.Disconnect(ctx)
}

// GetBucket returns a mongo collection based on the input collection name
func (db *dbClient) GetBucket(bucketName string) (*gridfs.Bucket, error) {
	bucketOpts := options.GridFSBucket()
	bucketOpts.SetName(bucketName)
	bucket, err := gridfs.NewBucket(db.client.Database(viper.GetString("database")), bucketOpts)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

// GetCollection returns a mongo collection based on the input collection name
func (db *dbClient) GetCollection(collectionName string) DBCollection {
	return db.client.Database(viper.GetString("database")).Collection(collectionName)
}

// NewDBHandler returns a new DBHandler generic interface
func (db *dbClient) NewDBHandler(collectionName string) *DBRepo[DBRecord] {
	col := db.GetCollection(collectionName)
	return &DBRepo[DBRecord]{
		db:         db,
		collection: col,
	}
}
