package repositories

import (
	"errors"
	"github.com/JECSand/eventit-server/domains/identity/src/models"
	"github.com/JECSand/eventit-server/domains/shared/databases"
	"github.com/JECSand/eventit-server/domains/shared/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// BlacklistRepo is used by the app to manage all user related controllers and functionality
type BlacklistRepo struct {
	collection       databases.DBCollection
	db               databases.DBClient
	blacklistHandler *databases.DBRepo[*BlacklistRecord]
}

// NewBlacklistRepo is an exported function used to initialize a new BlacklistRepo struct
func NewBlacklistRepo(db databases.DBClient) *BlacklistRepo {
	collection := db.GetCollection("users")
	repoHandler := &databases.DBRepo[*BlacklistRecord]{
		DB:         db,
		Collection: collection,
	}
	return &BlacklistRepo{collection, db, repoHandler}
}

type BlacklistRecord struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AuthToken string             `json:"auth_token" bson:"auth_token,omitempty"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
}

// newBlacklistRecord initializes a new pointer to a blacklistModel struct from a pointer to a JSON Blacklist struct
func newBlacklistRecord(bl *models.Blacklist) (bm *BlacklistRecord, err error) {
	bm = &BlacklistRecord{
		AuthToken: bl.AuthToken,
		UpdatedAt: bl.UpdatedAt,
		CreatedAt: bl.CreatedAt,
	}
	if bl.Id != "" && bl.Id != "000000000000000000000000" {
		bm.Id, err = primitive.ObjectIDFromHex(bl.Id)
	}
	return
}

// Update the blacklistModel using an overwrite bson doc
func (b *BlacklistRecord) Update(doc interface{}) (err error) {
	data, err := utilities.BSONMarshall(doc)
	if err != nil {
		return
	}
	bm := BlacklistRecord{}
	err = bson.Unmarshal(data, &bm)
	if len(bm.AuthToken) > 0 {
		b.AuthToken = bm.AuthToken
	}
	if !bm.UpdatedAt.IsZero() {
		b.UpdatedAt = bm.UpdatedAt
	}
	return
}

// BsonLoad loads a bson doc into the blacklistModel
func (b *BlacklistRecord) BsonLoad(doc bson.D) (err error) {
	bData, err := utilities.BSONMarshall(doc)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(bData, b)
	return err
}

// Match compares an input bson doc and returns whether there's a match with the blacklistModel
func (b *BlacklistRecord) Match(doc interface{}) bool {
	data, err := utilities.BSONMarshall(doc)
	if err != nil {
		return false
	}
	bm := BlacklistRecord{}
	err = bson.Unmarshal(data, &bm)
	if b.Id == bm.Id {
		return true
	}
	if b.AuthToken == bm.AuthToken {
		return true
	}
	return false
}

// GetID returns the unique identifier of the blacklistModel
func (b *BlacklistRecord) GetID() (id interface{}) {
	return b.Id
}

// AddTimeStamps updates a blacklistModel struct with a timestamp
func (b *BlacklistRecord) AddTimeStamps(newRecord bool) {
	currentTime := time.Now().UTC()
	b.UpdatedAt = currentTime
	if newRecord {
		b.CreatedAt = currentTime
	}
}

// AddObjectID checks if a blacklistModel has a value assigned for Id, if no value a new one is generated and assigned
func (b *BlacklistRecord) AddObjectID() {
	if b.Id.Hex() == "" || b.Id.Hex() == "000000000000000000000000" {
		b.Id = primitive.NewObjectID()
	}
}

// PostProcess updates an blacklistModel struct postProcess to do things such as removing the password field's value
func (b *BlacklistRecord) PostProcess() (err error) {
	if b.AuthToken == "" {
		err = errors.New("blacklist record does not have an AuthToken")
	}
	return
}

// ToDoc converts the bson blacklistModel into a bson.D
func (b *BlacklistRecord) ToDoc() (doc bson.D, err error) {
	data, err := bson.Marshal(b)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

// BsonFilter generates a bson filter for MongoDB queries from the blacklistModel data
func (b *BlacklistRecord) BsonFilter() (doc bson.D, err error) {
	if b.AuthToken != "" {
		doc = bson.D{{"auth_token", b.AuthToken}}
	} else if b.Id.Hex() != "" && b.Id.Hex() != "000000000000000000000000" {
		doc = bson.D{{"_id", b.Id}}
	}
	return
}

// BsonUpdate generates a bson update for MongoDB queries from the blacklistModel data
func (b *BlacklistRecord) BsonUpdate() (doc bson.D, err error) {
	inner, err := b.ToDoc()
	if err != nil {
		return
	}
	doc = bson.D{{"$set", inner}}
	return
}

// ToRoot creates and return a new pointer to a Blacklist JSON struct from a pointer to a BSON blacklistModel
func (b *BlacklistRecord) ToRoot() *models.Blacklist {
	return &models.Blacklist{
		Id:        b.Id.Hex(),
		AuthToken: b.AuthToken,
		UpdatedAt: b.UpdatedAt,
		CreatedAt: b.CreatedAt,
	}
}
