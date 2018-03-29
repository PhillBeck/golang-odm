package odm

import (
	"gopkg.in/mgo.v2/bson"
)

// DocumentBase implements all the methods to the IEntity interface
type DocumentBase struct {
	ID bson.ObjectId `json:"_id" bson:"_id"`
}

// SetID sets the ID field of the document
func (doc *DocumentBase) SetID(ID bson.ObjectId) {
	doc.ID = ID
}

// GetID returns the id of the document
func (doc *DocumentBase) GetID() bson.ObjectId {
	return doc.ID
}

// IEntity defines all the methods to be used on the
// documents when saved and queried on the db
type IEntity interface {
	SetID(bson.ObjectId)
	GetID() bson.ObjectId
}
