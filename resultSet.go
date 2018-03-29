package odm

import (
	"encoding/json"

	"gopkg.in/mgo.v2"
)

type Query struct {
	counter int
	query   *mgo.Query
	fetched bool
	results []interface{}
	session *mgo.Session
}

func (q *Query) Limit(n int) *Query {
	query := q.query.Limit(n)
	return &Query{
		counter: q.counter,
		fetched: q.fetched,
		session: q.session,
		query:   query}
}

func (q *Query) Skip(n int) *Query {
	query := q.query.Skip(n)
	return &Query{
		counter: q.counter,
		fetched: q.fetched,
		session: q.session,
		query:   query}
}

func (q *Query) Sort(fields ...string) *Query {
	query := q.query.Sort(fields...)
	return &Query{
		counter: q.counter,
		fetched: q.fetched,
		session: q.session,
		query:   query}
}

func (q *Query) All(docs interface{}) error {
	defer q.session.Close()
	return q.query.All(docs)
}

func (q *Query) Next(doc IEntity) bool {
	if !q.fetched {
		q.results = []interface{}{}
		if q.query.All(&q.results) != nil {
			return false
		}
		q.fetched = true
	}

	if q.counter+1 > len(q.results) {
		doc = nil
		return false
	}

	marshalled, err := json.Marshal(q.results[q.counter])
	if err != nil {
		return false
	}
	err = json.Unmarshal(marshalled, doc)
	if err != nil {
		return false
	}

	q.counter++

	return true
}
