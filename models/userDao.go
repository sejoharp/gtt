package models

import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"

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

func (dao *UserDao) FindByID(id bson.ObjectId) (User, error) {
	var user User
	err := dao.getDBCollection().FindId(id).One(&user)
	return user, err
}

func (dao *UserDao) FindByName(name string) (User, error) {
	var user User
	err := dao.getDBCollection().Find(bson.M{"name": name}).One(&user)
	return user, err
}
