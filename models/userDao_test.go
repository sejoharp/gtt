package models

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("UserDao", func() {
	var (
		collection *mgo.Collection
		dao        *UserDao
		id         bson.ObjectId
		name       string
		worktime   time.Duration
		overtime   time.Duration
	)

	BeforeEach(func() {
		session, err := createSession()
		Expect(err).To(BeNil(), "All tests need a connection to a mongodb.")
		collection = getCollection(session, "timetracker", "users")
		cleanCollection(collection)

		dao = NewUserDao(session, "timetracker")

		id = bson.NewObjectId()
		name = "myuser"
		worktime, _ = time.ParseDuration("7h45m")
		overtime, _ = time.ParseDuration("1h")
	})

	It("should save a user.", func() {
		user := NewPersistedUser(id, name, worktime, overtime)

		err := dao.Save(user)

		Expect(err).To(Succeed())
		var persistedUser User
		findErr := collection.FindId(id).One(&persistedUser)
		Expect(findErr).To(Succeed())
		checkPersistedUserEquality(user, persistedUser)
	})

	PIt("should find a user by id.")
	PIt("should find a user by name.")
	PIt("should return the password hash of a user.")
	PIt("should save a password hash of a user.")
})

func checkUserEquality(user1, user2 User) {
	Expect(user1.Name).To(Equal(user2.Name))
	Expect(user1.Overtime).To(Equal(user2.Overtime))
	Expect(user1.Worktime).To(Equal(user2.Worktime))
}

func checkPersistedUserEquality(user1, user2 User) {
	Expect(user1.ID).To(Equal(user2.ID))
	Expect(user1.Name).To(Equal(user2.Name))
	Expect(user1.Overtime).To(Equal(user2.Overtime))
	Expect(user1.Worktime).To(Equal(user2.Worktime))
}
