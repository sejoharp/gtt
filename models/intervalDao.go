package models

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type IntervalDao interface {
	Save(interval Interval) error
	FindByUserID(userID bson.ObjectId) ([]Interval, error)
	IsUserWorking(userID bson.ObjectId) (bool, error)
	Start(userID bson.ObjectId) error
	Stop(userID bson.ObjectId) error
	FindOpenIntervals(userID bson.ObjectId) ([]Interval, error)
	FindInRange(userID bson.ObjectId, begin time.Time, end time.Time) ([]Interval, error)
}

type IntervalDaoImpl struct {
	dbName         string
	collectionName string
	session        *mgo.Session
}

func NewIntervalDao(session *mgo.Session, dbName string) IntervalDao {
	return &IntervalDaoImpl{session: session, dbName: dbName, collectionName: "intervals"}
}

func (dao *IntervalDaoImpl) Save(interval Interval) error {
	return dao.getDBCollection().Insert(interval)
}

func (dao *IntervalDaoImpl) FindByUserID(userID bson.ObjectId) ([]Interval, error) {
	var intervals []Interval
	err := dao.getDBCollection().Find(bson.M{"userid": userID}).All(&intervals)
	return intervals, err
}

func (dao *IntervalDaoImpl) IsUserWorking(userID bson.ObjectId) (bool, error) {
	query := bson.M{"userid": userID, "stop": bson.M{"$exists": false}}
	openIntervals, err := dao.getDBCollection().Find(query).Count()
	working := openIntervals > 0
	return working, err
}

func (dao *IntervalDaoImpl) Start(userID bson.ObjectId) error {
	return dao.Save(NewIntervalWithStart(userID, time.Now()))
}

// Stop sets stop to current time, if the user is working. Returns an error when the user is not working or has more than one interval without a stop time or find returns an error.
func (dao *IntervalDaoImpl) Stop(userID bson.ObjectId) error {
	openIntervals, err := dao.FindOpenIntervals(userID)
	if validationErr := checkStopErrors(openIntervals, err); validationErr != nil {
		return validationErr
	}
	change := bson.M{"$set": bson.M{"stop": cleanNanoSeconds(time.Now())}}
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

// FindOpenIntervals returns all intervals where stop is not set/zero value
func (dao *IntervalDaoImpl) FindOpenIntervals(userID bson.ObjectId) ([]Interval, error) {
	var openIntervals []Interval
	query := bson.M{"userid": userID, "stop": bson.M{"$exists": false}}
	err := dao.getDBCollection().Find(query).All(&openIntervals)
	return openIntervals, err
}

// FindInRange returns all intervals where start is greater than or equals begin and less than end
func (dao *IntervalDaoImpl) FindInRange(userID bson.ObjectId, begin time.Time, end time.Time) ([]Interval, error) {
	var intervals []Interval
	query := bson.M{"userid": userID, "start": bson.M{"$gte": begin, "$lt": end}}
	err := dao.getDBCollection().Find(query).All(&intervals)
	return intervals, err
}

func (dao *IntervalDaoImpl) getDBConnection() *mgo.Database {
	return dao.session.Clone().DB(dao.dbName)
}

func (dao *IntervalDaoImpl) getDBCollection() *mgo.Collection {
	return dao.getDBConnection().C(dao.collectionName)
}
