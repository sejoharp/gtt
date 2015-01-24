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
		start := time.Now()
		stop := time.Now()

		Expect(dao.Save(NewInterval(userID, start, stop))).To(Succeed())

		var interval Interval
		findErr := collection.Find(bson.M{"userid": userID}).One(&interval)
		Expect(findErr).To(BeNil())
		Expect(interval.Start.Unix()).To(Equal(start.Unix()))
		Expect(interval.Stop.Unix()).To(Equal(stop.Unix()))
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

	It("should find a not working user.", func() {
		userID := bson.NewObjectId()
		dao.Save(NewInterval(userID, time.Now(), time.Now()))

		working, err := dao.IsUserWorking(userID)

		Expect(err).To(BeNil())
		Expect(working).To(BeFalse())
	})

	It("should find a working user.", func() {
		userID := bson.NewObjectId()
		dao.Save(NewIntervalWithStart(userID, time.Now()))

		working, err := dao.IsUserWorking(userID)

		Expect(err).To(BeNil())
		Expect(working).To(BeTrue())
	})

	It("should start a new interval with current start time.", func() {
		userID := bson.NewObjectId()

		Expect(dao.Start(userID)).To(Succeed())

		intervals, findErr := dao.FindByUserID(userID)
		Expect(findErr).To(BeNil())
		Expect(intervals[0].Stop).To(BeZero())
		Expect(intervals[0].Start).NotTo(BeZero())
	})

	It("should return all open intervals", func() {
		userID := bson.NewObjectId()
		dao.Start(userID)
		dao.Start(userID)
		dao.Save(NewInterval(userID, time.Now(), time.Now()))

		openIntervals, err := dao.FindOpenIntervals(userID)

		Expect(err).To(BeNil())
		Expect(openIntervals).To(HaveLen(2))
		Expect(openIntervals[0].Stop).To(BeZero())
		Expect(openIntervals[1].Stop).To(BeZero())
	})

	It("should return an error when stop does not find an open interval.", func() {
		userID := bson.NewObjectId()

		err := dao.Stop(userID)

		Expect(err.Error()).To(Equal("user is not working"))
	})

	It("should return an error when stop finds more than one open interval.", func() {
		userID := bson.NewObjectId()
		dao.Start(userID)
		dao.Start(userID)

		err := dao.Stop(userID)

		Expect(err.Error()).To(Equal("more than one open interval"))
	})

	It("should stop the last open interval.", func() {
		userID := bson.NewObjectId()
		dao.Start(userID)

		Expect(dao.Stop(userID)).To(Succeed())
	})

	PIt("should return all intervals in a given range.")
})
