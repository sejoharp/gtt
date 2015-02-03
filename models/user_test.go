package models

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("User", func() {
	id1 := bson.NewObjectId()
	id2 := bson.NewObjectId()
	name := "myuser"
	worktime, _ := time.ParseDuration("7h45m")
	overtime1, _ := time.ParseDuration("1h")
	overtime2, _ := time.ParseDuration("3h")

	user1 := NewUser(name, worktime, overtime1)
	user2 := NewUser(name, worktime, overtime2)
	user1Persisted := NewPersistedUser(id1, name, worktime, overtime1)
	user2Persisted := NewPersistedUser(id2, name, worktime, overtime2)

	It("should detect equal users.", func() {
		Expect(user1.EqualsWithoutID(user1)).To(BeTrue())
	})

	It("should detect different overtimes.", func() {
		Expect(user1.EqualsWithoutID(user2)).To(BeFalse())
	})

	It("should detect equal users.", func() {
		Expect(user1Persisted).To(Equal(user1Persisted))
	})

	It("should detect different ids.", func() {
		Expect(user1Persisted).NotTo(Equal(user2Persisted))
	})
})
