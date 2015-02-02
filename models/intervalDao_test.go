package models

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("IntervalDao", func() {
	var (
		collection *mgo.Collection
		dao        *IntervalDao
		userID bson.ObjectId
	)

	BeforeEach(func() {
		session, err := createSession()
		Expect(err).To(BeNil(), "All tests need a connection to a mongodb.")
		collection = getCollection(session, "timetracker", "intervals")
		cleanCollection(collection)

		dao = NewIntervalDao(session, "timetracker")

		userID = bson.NewObjectId()
	})

	It("should save an interval.", func() {
		expectedInterval := NewInterval(userID, time.Now(), time.Now())

		Expect(dao.Save(expectedInterval)).To(Succeed())

		var interval Interval
		findErr := collection.Find(bson.M{"userid": userID}).One(&interval)
		Expect(findErr).To(BeNil())
		Expect(interval.EqualsWithoutID(expectedInterval)).To(BeTrue())
		Expect(interval.ID.Valid()).To(BeTrue())
	})
	It("should save a complete interval.", func() {
		expectedInterval := NewPersistedInterval(bson.NewObjectId(), userID, time.Now(), time.Now())

		Expect(dao.Save(expectedInterval)).To(Succeed())

		var interval Interval
		Expect(collection.Find(bson.M{"userid": userID}).One(&interval)).To(Succeed())
		Expect(interval).To(Equal(expectedInterval))
	})

	It("should find all by userID.", func() {
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
		Expect(intervals).To(BeEmpty())
	})

	It("should find a not working user.", func() {
		dao.Save(NewInterval(userID, time.Now(), time.Now()))

		working, err := dao.IsUserWorking(userID)

		Expect(err).To(BeNil())
		Expect(working).To(BeFalse())
	})

	It("should find a working user.", func() {
		dao.Save(NewIntervalWithStart(userID, time.Now()))

		working, err := dao.IsUserWorking(userID)

		Expect(err).To(BeNil())
		Expect(working).To(BeTrue())
	})

	It("should start a new interval with current start time.", func() {
		Expect(dao.Start(userID)).To(Succeed())

		intervals, findErr := dao.FindByUserID(userID)
		Expect(findErr).To(BeNil())
		Expect(intervals[0].Stop).To(BeZero())
		Expect(intervals[0].Start).NotTo(BeZero())
	})

	It("should return all open intervals", func() {
		dao.Start(userID)
		dao.Start(userID)
		dao.Save(NewInterval(userID, time.Now(), time.Now()))

		openIntervals, err := dao.FindOpenIntervals(userID)

		Expect(err).To(BeNil())
		Expect(openIntervals).To(HaveLen(2))
		Expect(openIntervals[0].Stop).To(BeZero())
		Expect(openIntervals[1].Stop).To(BeZero())
	})

	It("should return one open interval", func() {
		dao.Start(userID)

		openIntervals, err := dao.FindOpenIntervals(userID)

		Expect(err).To(BeNil())
		Expect(openIntervals).To(HaveLen(1))
		Expect(openIntervals[0].Start).NotTo(BeZero())
		Expect(openIntervals[0].Stop).To(BeZero())
	})

	It("should return an error when stop does not find an open interval.", func() {
		err := dao.Stop(userID)

		Expect(err.Error()).To(Equal("user is not working"))
	})

	It("should return an error when stop finds more than one open interval.", func() {
		dao.Start(userID)
		dao.Start(userID)

		err := dao.Stop(userID)

		Expect(err.Error()).To(Equal("more than one open interval"))
	})

	It("should stop open interval.", func() {
		dao.Start(userID)

		Expect(dao.Stop(userID)).To(Succeed())

		intervals, err := dao.FindByUserID(userID)
		Expect(err).To(BeNil())
		Expect(intervals).To(HaveLen(1))
		Expect(intervals[0].Stop).ToNot(BeZero())
	})

	It("should return all intervals, where start is in a given range.", func() {
		interval := NewInterval(userID, createDate("2014-12-10 07:00"), createDate("2014-12-10 15:00"))
		dao.Save(NewInterval(userID, createDate("2014-12-11 08:00"), createDate("2014-12-11 16:00")))
		dao.Save(interval)
		dao.Save(NewInterval(userID, createDate("2014-12-09 09:00"), createDate("2014-12-09 19:00")))

		intervalsInRange, err := dao.FindInRange(userID, createDate("2014-12-10 00:00"), createDate("2014-12-11 00:00"))

		Expect(err).To(BeNil())
		Expect(intervalsInRange).To(HaveLen(1))
		Expect(intervalsInRange[0].EqualsWithoutID(interval)).To(BeTrue())
	})

	It("should return an empty array when there are no intervals in range.", func() {
		dao.Save(NewInterval(userID, createDate("2014-12-09 08:00"), createDate("2014-12-09 23:59")))
		dao.Save(NewInterval(userID, createDate("2014-12-11 00:00"), createDate("2014-12-11 16:00")))

		intervalsInRange, err := dao.FindInRange(userID, createDate("2014-12-10 00:00"), createDate("2014-12-11 00:00"))

		Expect(err).To(BeNil())
		Expect(intervalsInRange).To(BeEmpty())
	})

	It("should return all intervals, where start is near the limits in a given range.", func() {
		interval1 := NewInterval(userID, time.Date(2014, 12, 9, 23, 59, 59, 999, time.Local), createDate("2014-12-10 07:00"))
		interval2 := NewInterval(userID, time.Date(2014, 12, 10, 0, 0, 0, 0, time.Local), createDate("2014-12-10 07:00"))
		interval3 := NewInterval(userID, time.Date(2014, 12, 10, 23, 59, 59, 999, time.Local), createDate("2014-12-11 07:00"))
		interval4 := NewInterval(userID, time.Date(2014, 12, 11, 0, 0, 0, 0, time.Local), createDate("2014-12-10 07:00"))
		Expect(collection.Insert(interval1, interval2, interval3, interval4)).To(Succeed())

		intervalsInRange, err := dao.FindInRange(userID, createDate("2014-12-10 00:00"), createDate("2014-12-11 00:00"))

		Expect(err).To(BeNil())
		Expect(intervalsInRange).To(HaveLen(2))
		Expect(intervalsInRange[0].EqualsWithoutID(interval2)).To(BeTrue())
		Expect(intervalsInRange[1].EqualsWithoutID(interval3)).To(BeTrue())
	})
})

func createDate(date string) time.Time {
	time, err := time.ParseInLocation("2006-01-02 15:04", date, time.Local)
	if err != nil {
		panic(fmt.Sprintf("date parsing failed|date: %s", date))
	}
	return time
}
