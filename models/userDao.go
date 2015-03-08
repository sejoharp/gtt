package models

import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"

type UserDao interface {
	Save(user User) error
	SaveWithPassword(user UserWithPassword) error
	FindByID(id bson.ObjectId) (User, error)
	FindByName(name string) (User, error)
	AddPassword(id bson.ObjectId, password string) error
	AddPasswordByUser(username string, password string) error
	GetPassword(id bson.ObjectId) (string, error)
	GetPasswordByUser(username string) (string, error)
	Update(user User) error
}

type UserDaoImpl struct {
	dbName         string
	collectionName string
	session        *mgo.Session
}

func NewUserDao(session *mgo.Session, dbName string) UserDao {
	return &UserDaoImpl{session: session, dbName: dbName, collectionName: "users"}
}

func (dao *UserDaoImpl) Save(user User) error {
	return dao.getDBCollection().Insert(user)
}

func (dao *UserDaoImpl) SaveWithPassword(user UserWithPassword) error {
	return dao.getDBCollection().Insert(user)
}

func (dao *UserDaoImpl) getDBConnection() *mgo.Database {
	return dao.session.Clone().DB(dao.dbName)
}

func (dao *UserDaoImpl) getDBCollection() *mgo.Collection {
	return dao.getDBConnection().C(dao.collectionName)
}

func (dao *UserDaoImpl) FindByID(id bson.ObjectId) (User, error) {
	var user User
	err := dao.getDBCollection().FindId(id).One(&user)
	return user, err
}

func (dao *UserDaoImpl) FindByName(name string) (User, error) {
	var user User
	err := dao.getDBCollection().Find(bson.M{"name": name}).One(&user)
	return user, err
}

func (dao *UserDaoImpl) AddPassword(id bson.ObjectId, password string) error {
	change := bson.M{"$set": bson.M{"password": password}}
	return dao.getDBCollection().UpdateId(id, change)
}

func (dao *UserDaoImpl) AddPasswordByUser(username string, password string) error {
	change := bson.M{"$set": bson.M{"password": password}}
	query := bson.M{"name": username}
	return dao.getDBCollection().Update(query, change)
}

func (dao *UserDaoImpl) GetPassword(id bson.ObjectId) (string, error) {
	var result bson.M
	err := dao.getDBCollection().FindId(id).Select(bson.M{"password": 1}).One(&result)
	return result["password"].(string), err
}

func (dao *UserDaoImpl) Update(user User) error {
	return dao.getDBCollection().UpdateId(user.ID, user)
}

func (dao *UserDaoImpl) GetPasswordByUser(username string) (string, error) {
	var result bson.M
	query := bson.M{"name": username}
	err := dao.getDBCollection().Find(query).Select(bson.M{"password": 1}).One(&result)
	return result["password"].(string), err
}
