package models

import (
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Interval contains the Duration of a work from a user
type Interval struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	UserID bson.ObjectId
	Start  time.Time
	Stop   time.Time `bson:",omitempty"`
}

type IntervalDao struct {
	dbName         string
	collectionName string
	session        *mgo.Session
}

func NewInterval(userID bson.ObjectId, start time.Time, stop time.Time) Interval {
	return Interval{UserID: userID, Start: start, Stop: stop}
}

func NewIntervalWithStart(userID bson.ObjectId, start time.Time) Interval {
	return Interval{UserID: userID, Start: start}
}

func NewPersistedInterval(id bson.ObjectId, userID bson.ObjectId, start time.Time, stop time.Time) Interval {
	return Interval{id, userID, start, stop}
}

func NewPersistedIntervalWithStart(userID bson.ObjectId, start time.Time) Interval {
	return Interval{UserID: userID, Start: start}
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
