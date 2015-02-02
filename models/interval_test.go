package models

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("Interval", func() {
	time1 := time.Date(2014, 12, 14, 8, 10, 11, 99, time.Local)
	time1NoNanoSeconds := time.Date(2014, 12, 14, 8, 10, 11, 0, time.Local)
	time2 := time.Date(2014, 12, 14, 8, 10, 12, 1, time.Local)

	userID := bson.NewObjectId()
	id := bson.NewObjectId()

	interval1 := NewPersistedInterval(id, userID, time1, time1)
	interval1NoNanoSeconds := NewPersistedInterval(id, userID, time1NoNanoSeconds, time1NoNanoSeconds)
	interval2 := NewPersistedInterval(id, userID, time2, time2)
	interval3 := NewInterval(userID, time2, time2)
	interval4 := NewIntervalWithStart(userID, time2)
	interval5 := NewInterval(bson.NewObjectId(), time2, time2)

	It("should detect equal intervals.", func() {
		Expect(interval1).To(Equal(interval1))
	})

	It("should detect different intervals.", func() {
		Expect(interval1).NotTo(Equal(interval2))
	})

	It("should detect different userIDs.", func() {
		Expect(interval3.EqualsWithoutID(interval5)).To(BeFalse())
	})

	It("should ignore different nanoseconds.", func() {
		Expect(interval1).To(Equal(interval1NoNanoSeconds))
	})

	It("should detect a missing id.", func() {
		Expect(interval3).NotTo(Equal(interval2))
	})

	It("should detect a missing stop date.", func() {
		Expect(interval3).NotTo(Equal(interval4))
	})

	It("should clean milli- and nanoseconds.", func() {
		expectedTime := time.Date(2014, 12, 14, 8, 10, 11, 0, time.Local)
		Expect(cleanNanoSeconds(time1)).To(Equal(expectedTime))
	})
})
