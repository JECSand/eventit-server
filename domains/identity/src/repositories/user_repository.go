package repositories

import (
	"errors"
	"github.com/JECSand/eventit-server/domains/identity/src/models"
	"github.com/JECSand/eventit-server/domains/shared/databases"
	"github.com/JECSand/eventit-server/domains/shared/enums"
	"github.com/JECSand/eventit-server/domains/shared/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// UserRepo is used by the app to manage all user related controllers and functionality
type UserRepo struct {
	collection databases.DBCollection
	db         databases.DBClient
	Handler    *databases.DBRepo[*UserRecord]
}

// NewUserRepo is an exported function used to initialize a new UserRepo struct
func NewUserRepo(db databases.DBClient) *UserRepo {
	collection := db.GetCollection("users")
	repoHandler := &databases.DBRepo[*UserRecord]{
		DB:         db,
		Collection: collection,
	}
	return &UserRepo{collection, db, repoHandler}
}

// UserRecord stores User information
type UserRecord struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username,omitempty"`
	Password  string             `json:"password" bson:"password,omitempty"`
	FirstName string             `json:"firstname" bson:"firstname,omitempty"`
	LastName  string             `json:"lastname" bson:"lastname,omitempty"`
	Email     string             `json:"email" bson:"email,omitempty"`
	Role      enums.Role         `json:"role" bson:"role,omitempty"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
	DeletedAt time.Time          `json:"deleted_at" bson:"deleted_at,omitempty"`
}

// NewUserRecord initializes a new pointer to a UserRecord struct from a pointer to a JSON User struct
func NewUserRecord(u *models.User) (um *UserRecord, err error) {
	um = &UserRecord{
		Username:  u.Username,
		Password:  u.Password,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Role:      u.Role,
		UpdatedAt: u.UpdatedAt,
		CreatedAt: u.CreatedAt,
		DeletedAt: u.DeletedAt,
	}
	if u.Id != "" && u.Id != "000000000000000000000000" {
		um.Id, err = primitive.ObjectIDFromHex(u.Id)
	}
	return
}

// Update the UserRecord using an overwrite bson.D doc
func (u *UserRecord) Update(doc interface{}) (err error) {
	data, err := utilities.BSONMarshall(doc)
	if err != nil {
		return
	}
	um := UserRecord{}
	err = bson.Unmarshal(data, &um)
	if len(um.Username) > 0 {
		u.Username = um.Username
	}
	if len(um.FirstName) > 0 {
		u.FirstName = um.FirstName
	}
	if len(um.LastName) > 0 {
		u.LastName = um.LastName
	}
	if len(um.Email) > 0 {
		u.Email = um.Email
	}
	if len(um.Password) > 0 {
		u.Password = um.Password
	}
	if um.Role.EnumIndex() > 0 {
		u.Role = um.Role
	}
	if !um.UpdatedAt.IsZero() {
		u.UpdatedAt = um.UpdatedAt
	}
	return
}

// BsonLoad loads a bson doc into the UserRecord
func (u *UserRecord) BsonLoad(doc bson.D) (err error) {
	bData, err := utilities.BSONMarshall(doc)
	if err != nil {
		return
	}
	err = bson.Unmarshal(bData, u)
	return
}

// Match compares an input bson doc and returns whether there's a match with the UserRecord
// TODO: Find a better way to write these model match methods
func (u *UserRecord) Match(doc interface{}) bool {
	data, err := utilities.BSONMarshall(doc)
	if err != nil {
		return false
	}
	um := UserRecord{}
	err = bson.Unmarshal(data, &um)
	if um.Id.Hex() != "" && um.Id.Hex() != "000000000000000000000000" {
		if u.Id == um.Id {
			return true
		}
		return false
	}
	if um.Email != "" {
		if u.Email == um.Email {
			return true
		}
		return false
	}
	return false
}

// GetID returns the unique identifier of the UserRecord
func (u *UserRecord) GetID() (id interface{}) {
	return u.Id
}

// AddTimeStamps updates an UserRecord struct with a timestamp
func (u *UserRecord) AddTimeStamps(newRecord bool) {
	currentTime := time.Now().UTC()
	u.UpdatedAt = currentTime
	if newRecord {
		u.CreatedAt = currentTime
	}
}

// AddObjectID checks if a userModel has a value assigned for Id if no value a new one is generated and assigned
func (u *UserRecord) AddObjectID() {
	if u.Id.Hex() == "" || u.Id.Hex() == "000000000000000000000000" {
		u.Id = primitive.NewObjectID()
	}
}

// PostProcess updates an userModel struct postProcess to do things such as removing the password field's value
func (u *UserRecord) PostProcess() (err error) {
	//u.Password = ""
	if u.Email == "" {
		err = errors.New("user record does not have an email")
	}
	// TODO - When implementing soft delete, DeletedAt can be checked here to ensure deleted users are filtered out
	return
}

// ToDoc converts the bson userModel into a bson.D
func (u *UserRecord) ToDoc() (doc bson.D, err error) {
	data, err := bson.Marshal(u)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

// BsonFilter generates a bson filter for MongoDB queries from the userModel data
func (u *UserRecord) BsonFilter() (doc bson.D, err error) {
	if u.Id.Hex() != "" && u.Id.Hex() != "000000000000000000000000" {
		doc = bson.D{{"_id", u.Id}}
	} else if u.Email != "" {
		doc = bson.D{{"email", u.Email}}
	}
	return
}

// BsonUpdate generates a bson update for MongoDB queries from the userModel data
func (u *UserRecord) BsonUpdate() (doc bson.D, err error) {
	inner, err := u.ToDoc()
	if err != nil {
		return
	}
	doc = bson.D{{"$set", inner}}
	return
}

// ToRoot creates and return a new pointer to a User JSON struct from a pointer to a BSON userModel
func (u *UserRecord) ToRoot() *models.User {
	return &models.User{
		Id:        u.Id.Hex(),
		Username:  u.Username,
		Password:  u.Password,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Role:      u.Role,
		UpdatedAt: u.UpdatedAt,
		CreatedAt: u.CreatedAt,
		DeletedAt: u.DeletedAt,
	}
}

// LoadUserRecords ..
func LoadUserRecords(ms []*UserRecord) (users []*models.User) {
	for _, m := range ms {
		users = append(users, m.ToRoot())
	}
	return
}
