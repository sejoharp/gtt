package models

import (
	"time"

	"errors"

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

func (dao *IntervalDao) IsUserWorking(userID bson.ObjectId) (bool, error) {
	query := bson.M{"userid": userID, "stop": bson.M{"$exists": false}}
	openIntervals, err := dao.getDBCollection().Find(query).Count()
	working := openIntervals > 0
	return working, err
}

func (dao *IntervalDao) Start(userID bson.ObjectId) error {
	return dao.Save(NewIntervalWithStart(userID, time.Now()))
}

func (dao *IntervalDao) Stop(userID bson.ObjectId) error {
	openIntervals, err := dao.FindOpenIntervals(userID)
	if validationErr := checkStopErrors(openIntervals, err); validationErr != nil {
		return validationErr
	}
	change := bson.M{"$set": bson.M{"stop": time.Now()}}
	return dao.getDBCollection().UpdateId(openIntervals[0].ID, change)
}

func checkStopErrors(openIntervals []Interval, err error) error {
	if err != nil {
		return err
	}
	if len(openIntervals) > 1 {
		return errors.New("more than one open interval")
	}
	if len(openIntervals) == 0 {
		return errors.New("user is not working")
	}
	return nil
}

func (dao *IntervalDao) FindOpenIntervals(userID bson.ObjectId) ([]Interval, error) {
	var openIntervals []Interval
	findQuery := bson.M{"userid": userID, "stop": bson.M{"$exists": false}}
	err := dao.getDBCollection().Find(findQuery).All(&openIntervals)
	return openIntervals, err
}

func (dao *IntervalDao) getDBConnection() *mgo.Database {
	return dao.session.Clone().DB(dao.dbName)
}

func (dao *IntervalDao) getDBCollection() *mgo.Collection {
	return dao.getDBConnection().C(dao.collectionName)
}
