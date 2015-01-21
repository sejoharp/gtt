package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2"

	"testing"
)

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}

func createSession() (*mgo.Session, error) {
	return mgo.Dial("localhost")
}

func getCollection(session *mgo.Session, dbName string, collectionName string) *mgo.Collection {
	return session.DB(dbName).C(collectionName)
}

func cleanCollection(collection *mgo.Collection) error {
	return collection.DropCollection()
}
