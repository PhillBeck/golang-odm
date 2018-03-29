package odm

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mgoSession *mgo.Session
var MongoHost, MongoPort string

func getSession() (*mgo.Session, error) {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(MongoHost + ":" + MongoPort)
		if err != nil {
			return nil, err
		}
	}

	return mgoSession.Clone(), nil
}

func NewRepository(database, collection string) *Repository {
	return &Repository{
		database:   database,
		collection: collection,
	}
}

type Repository struct {
	database   string
	collection string
}

func (r *Repository) Save(doc IEntity) error {
	session, err := getSession()
	if err != nil {
		return err
	}
	defer session.Close()

	if doc.GetID().Hex() == "" {
		doc.SetID(bson.NewObjectId())
	}

	_, err = session.DB(r.database).C(r.collection).UpsertId(doc.GetID(), doc)

	return err
}

func (r *Repository) GetByID(ID bson.ObjectId, doc IEntity) error {
	session, err := getSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.DB(r.database).C(r.collection).FindId(ID).One(doc)
}

func (r *Repository) Delete(doc IEntity) error {
	id := doc.GetID()
	if id.Hex() == "" {
		return fmt.Errorf("Invalid ID")
	}

	session, err := getSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.DB(r.database).C(r.collection).RemoveId(doc.GetID())
}

func (r *Repository) Find(query bson.M) (*Query, error) {
	session, err := getSession()
	if err != nil {
		return nil, err
	}

	ret := Query{
		query:   session.DB(r.database).C(r.collection).Find(query),
		session: session,
		fetched: false}

	return &ret, nil
}

func (r *Repository) FindOne(query bson.M, doc IEntity) error {
	session, err := getSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.DB(r.database).C(r.collection).Find(query).One(doc)
}
