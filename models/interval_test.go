package models

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type IntervalDaoSuite struct {
	collection *mgo.Collection
	dao        *IntervalDao
}

func (suite *IntervalDaoSuite) SetUpTest(c *C) {
	suite.dao, _ = createDao()
	suite.collection = suite.dao.getDBCollection()
	cleanCollection(suite.collection)
}

func (suite *IntervalDaoSuite) TestSave(c *C) {
	userID := bson.NewObjectId()
	now := time.Now()

	insertErr := suite.dao.Save(NewIntervalStart(userID, now))

	c.Assert(insertErr, IsNil)
	var interval Interval
	findErr := suite.collection.Find(bson.M{"userid": userID}).One(&interval)
	c.Assert(findErr, IsNil)
	c.Assert(interval.Start.Unix(), Equals, now.Unix())
	c.Assert(interval.UserID, Equals, userID)
}

func (suite *IntervalDaoSuite) TestFindByUserID(c *C) {
	userID := bson.NewObjectId()
	suite.dao.Save(NewIntervalStart(userID, time.Now()))
	suite.dao.Save(NewIntervalStart(userID, time.Now()))
	suite.dao.Save(NewIntervalStart(bson.NewObjectId(), time.Now()))

	intervals, err := suite.dao.FindByUserID(userID)

	c.Assert(err, IsNil)
	c.Assert(intervals, HasLen, 2)
	c.Assert(intervals[0].UserID, Equals, userID)
	c.Assert(intervals[1].UserID, Equals, userID)
}

func createDao() (*IntervalDao, error) {
	session, err := createSession()
	dao := NewIntervalDao(session, "timetracker")
	return dao, err
}

func createSession() (*mgo.Session, error) {
	return mgo.Dial("localhost")
}

func cleanCollection(collection *mgo.Collection) error {
	return collection.DropCollection()
}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&IntervalDaoSuite{})
