package utilities

import (
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
