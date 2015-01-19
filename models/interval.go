package models

import (
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Interval contains the Duration of a work from a user
type Interval struct {
	UserID bson.ObjectId
	Start  time.Time
	Stop   time.Time `bson:",omitempty"`
}

type IntervalDao struct {
	dbName         string
	collectionName string
	session        *mgo.Session
}

func NewIntervalStart(userID bson.ObjectId, start time.Time) Interval {
	return Interval{UserID: userID, Start: start}
}

func NewInterval(userID bson.ObjectId, start time.Time, stop time.Time) Interval {
	return Interval{userID, start, stop}
}

func NewIntervalDao(session *mgo.Session, dbName string) *IntervalDao {
	return &IntervalDao{session: session, dbName: dbName, collectionName: "intervals"}
}

func (dao *IntervalDao) Save(interval Interval) error {
	return dao.getDBCollection().Insert(interval)
}

func (dao *IntervalDao) FindByUserID(userID bson.ObjectId) ([]Interval, error) {
	var intervals []Interval
	err := dao.getDBCollection().Find(bson.M{"userid": userID}).All(&intervals)
	return intervals, err
}

func (dao *IntervalDao) getDBConnection() *mgo.Database {
	return dao.session.Clone().DB(dao.dbName)
}

func (dao *IntervalDao) getDBCollection() *mgo.Collection {
	return dao.getDBConnection().C(dao.collectionName)
}
