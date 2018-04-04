package odm

import (
	"fmt"
	"math"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mgoSession *mgo.Session
var MongodbURI string

func getSession() (*mgo.Session, error) {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(MongodbURI)
		if err != nil {
			fmt.Printf("URI: %s\n", MongodbURI)
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

func (r *Repository) Paginate(query bson.M, recordsPerPage, page int) (*Query, *PaginationInfo, error) {
	if page < 1 {
		return nil, nil, fmt.Errorf("Invalid page. Should be > 1")
	}

	if recordsPerPage < 1 {
		return nil, nil, fmt.Errorf("Invalid recordsPerPage. Should be > 1")
	}

	session, err := getSession()
	if err != nil {
		return nil, nil, err
	}
	defer session.Close()

	col := session.DB(r.database).C(r.collection)

	numRecords, err := col.Find(query).Count()
	if err != nil {
		return nil, nil, err
	}

	skip := recordsPerPage * (page - 1)

	info := PaginationInfo{
		CurrentPage:    page,
		NumPages:       int(math.Ceil(float64(numRecords) / float64(recordsPerPage))),
		RecordsPerPage: recordsPerPage,
		NumRecords:     numRecords}

	resultSet, err := r.Find(query)
	if err != nil {
		return nil, nil, err
	}

	resultSet = resultSet.Skip(skip).Limit(recordsPerPage)

	return resultSet, &info, nil
}

func (r *Repository) Count(query bson.M) (int, error) {
	session, err := getSession()
	if err != nil {
		return 0, err
	}

	defer session.Close()

	return session.DB(r.database).C(r.collection).Find(query).Count()
}
