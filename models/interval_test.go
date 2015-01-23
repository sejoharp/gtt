package models_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/zippelmann/gtt/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("IntervalDao", func() {
	var (
		collection *mgo.Collection
		dao        *IntervalDao
	)

	BeforeEach(func() {
		session, err := createSession()
		Expect(err).To(BeNil(), "All tests need a connection to a mongodb.")

		collection = getCollection(session, "timetracker", "intervals")

		cleanCollection(collection)

		dao = NewIntervalDao(session, "timetracker")
	})

	It("should save an interval.", func() {
		userID := bson.NewObjectId()
		now := time.Now()

		insertErr := dao.Save(NewIntervalWithStart(userID, now))

		Expect(insertErr).To(BeNil())
		var interval Interval
		findErr := collection.Find(bson.M{"userid": userID}).One(&interval)
		Expect(findErr).To(BeNil())
		Expect(interval.Start.Unix()).To(Equal(now.Unix()))
		Expect(interval.UserID).To(Equal(userID))
		Expect(interval.ID.Valid()).To(BeTrue())
	})

	It("should find all by userID.", func() {
		userID := bson.NewObjectId()
		dao.Save(NewIntervalWithStart(userID, time.Now()))
		dao.Save(NewIntervalWithStart(userID, time.Now()))
		dao.Save(NewIntervalWithStart(bson.NewObjectId(), time.Now()))

		intervals, err := dao.FindByUserID(userID)

		Expect(err).To(BeNil())
		Expect(intervals).To(HaveLen(2))
		Expect(intervals[0].UserID).To(Equal(userID))
		Expect(intervals[1].UserID).To(Equal(userID))
	})

	It("should return an empty array when no intervals found.", func() {
		intervals, err := dao.FindByUserID(bson.NewObjectId())

		Expect(err).To(BeNil())
		Expect(intervals).To(HaveLen(0))
	})

	PIt("should find a not working user.")
	PIt("should find a working user.")
	PIt("should return all intervals in a given range.")
	PIt("should start a new interval.")
	PIt("should stop the last open interval.")
})

var _ = Describe("Interval", func() {
	It("should ensure a present userid", func() {
		var id bson.ObjectId = "a"
		Expect(id.Valid()).To(BeFalse())
	})

	PIt("should ensure a present start")
})
