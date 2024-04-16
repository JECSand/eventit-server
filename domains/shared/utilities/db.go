package utilities

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GenerateObjectID for index keying records of data
func GenerateObjectID() string {
	newId := primitive.NewObjectID()
	return newId.Hex()
}

// CheckObjectID checks whether a hexID is null or now
func CheckObjectID(hexID string) bool {
	if hexID == "" || hexID == "000000000000000000000000" {
		return false
	}
	return true
}

// BSONMarshall inputs a bson type and attempts to marshall it into a slice of bytes
func BSONMarshall(bsonData interface{}) (data []byte, err error) {
	switch t := bsonData.(type) {
	case nil:
		return nil, errors.New("input bsonData to marshall can not be nil")
	case []byte:
		return t, nil
	case bson.D:
		data, err = bson.Marshal(t)
		if err != nil {
			return
		}
	case bson.M:
		data, err = bson.Marshal(t)
		if err != nil {
			return
		}
	case bson.Raw:
		data, err = bson.Marshal(t)
		if err != nil {
			return
		}
	case bson.E:
		data, err = bson.Marshal(t)
		if err != nil {
			return
		}
	}
	return
}
