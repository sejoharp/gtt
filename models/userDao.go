package models

import "gopkg.in/mgo.v2"

type UserDao struct {
	dbName         string
	collectionName string
	session        *mgo.Session
}

func NewUserDao(session *mgo.Session, dbName string) *UserDao {
	return &UserDao{session: session, dbName: dbName, collectionName: "users"}
}

//TODO: function that return the password hash from a user

func (dao *UserDao) Save(user User) error {
	return dao.getDBCollection().Insert(user)
}

func (dao *UserDao) getDBConnection() *mgo.Database {
	return dao.session.Clone().DB(dao.dbName)
}

func (dao *UserDao) getDBCollection() *mgo.Collection {
	return dao.getDBConnection().C(dao.collectionName)
}
