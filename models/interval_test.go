package models_test

import (
	"time"

	. "github.com/zippelmann/gtt/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("IntervalDao", func() {
	var (
		collection *mgo.Collection
		dao        *IntervalDao
	)

	BeforeEach(func() {
		dao, _ = createDao()
		collection = dao.GetDBCollection()
		cleanCollection(collection)
	})

	It("should save an interval.", func() {
		userID := bson.NewObjectId()
		now := time.Now()

		insertErr := dao.Save(NewIntervalStart(userID, now))

		Expect(insertErr).To(BeNil())
		var interval Interval
		findErr := collection.Find(bson.M{"userid": userID}).One(&interval)
		Expect(findErr).To(BeNil())
		Expect(interval.Start.Unix()).To(Equal(now.Unix()))
		Expect(interval.UserID).To(Equal(userID))
	})
	It("should find all by userID.", func() {
		userID := bson.NewObjectId()
		dao.Save(NewIntervalStart(userID, time.Now()))
		dao.Save(NewIntervalStart(userID, time.Now()))
		dao.Save(NewIntervalStart(bson.NewObjectId(), time.Now()))

		intervals, err := dao.FindByUserID(userID)

		Expect(err).To(BeNil())
		Expect(intervals).To(HaveLen(2))
		Expect(intervals[0].UserID).To(Equal(userID))
		Expect(intervals[1].UserID).To(Equal(userID))
	})

	It("should return an empty array when no intervals found", func() {
		intervals, err := dao.FindByUserID(bson.NewObjectId())

		Expect(err).To(BeNil())
		Expect(intervals).To(HaveLen(0))
	})
})

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
