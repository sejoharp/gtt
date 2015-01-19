package models

import (
	"testing"
	"time"

	. "github.com/franela/goblin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Test(t *testing.T) {
	var collection *mgo.Collection
	var dao *IntervalDao

	g := Goblin(t)
	g.Describe("IntervalDao", func() {

		g.BeforeEach(func() {
			var err error
			dao, err = createDao()
			g.Assert(err).Equal(nil)
			collection = dao.getDBCollection()
			cleanCollection(collection)
		})

		g.It("should save an Interval", func() {
			userID := bson.NewObjectId()
			now := time.Now()

			insertErr := dao.Save(NewIntervalStart(userID, now))

			g.Assert(insertErr).Equal(nil)
			var interval Interval
			findErr := collection.Find(bson.M{"userid": userID}).One(&interval)
			g.Assert(findErr).Equal(nil)
			g.Assert(interval.Start.Unix()).Equal(now.Unix())
			g.Assert(interval.UserID).Equal(userID)
		})

		g.It("should find all intervals by userID", func() {
			userID := bson.NewObjectId()
			dao.Save(NewIntervalStart(userID, time.Now()))
			dao.Save(NewIntervalStart(userID, time.Now()))
			dao.Save(NewIntervalStart(bson.NewObjectId(), time.Now()))

			intervals, err := dao.FindByUserID(userID)

			g.Assert(err).Equal(nil)
			g.Assert(len(intervals)).Equal(2)
			g.Assert(intervals[0].UserID).Equal(userID)
			g.Assert(intervals[1].UserID).Equal(userID)
		})
	})
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
