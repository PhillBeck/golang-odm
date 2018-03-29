package odm_test

import (
	"reflect"
	"testing"

	odm "github.com/PhillBeck/golang-odm"
	"gopkg.in/mgo.v2/bson"
)

var repo *odm.Repository

type TestType struct {
	odm.DocumentBase `bson:",inline"`
	Name             string
	Description      string
}

func TestRepo(t *testing.T) {
	odm.MongoHost = "localhost"
	odm.MongoPort = "27000"
	repo = odm.NewRepository("odm-test", "test")

	for i := 0; i < 5; i++ {
		err := saveTest(createTest())
		if err != nil {
			t.Errorf("%+v\n", err)
		}
	}

	resultSet, err := repo.Find(bson.M{})
	if err != nil {
		t.Errorf("%+v\n", err)
	}

	var result []TestType

	err = resultSet.All(&result)
	if err != nil {
		t.Errorf("%+v\n", err)
	}

	if len(result) < 5 {
		t.Errorf("Not fetched")
	}

	resultSet1, err := repo.Find(bson.M{})
	if err != nil {
		t.Errorf("%+v\n", err)
	}

	doc := TestType{}
	for resultSet1.Next(&doc) {
		if reflect.TypeOf(doc) != reflect.TypeOf(TestType{}) {
			t.Errorf("Wrong types")
		}

		err = repo.Delete(&doc)
		if err != nil {
			t.Errorf("%+v\n", err)
		}
	}
}

func createTest() *TestType {
	return &TestType{
		Name:        "testName",
		Description: "testDescription"}
}

func saveTest(doc *TestType) error {
	return repo.Save(doc)
}
