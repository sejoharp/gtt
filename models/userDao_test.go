package models

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("UserDao", func() {
	const dbname = "timetracker"
	const collectionname = "users"

	var (
		collection  *mgo.Collection
		dao         *UserDao
		id          bson.ObjectId
		name        string
		worktime    time.Duration
		overtime    time.Duration
		userWithID  User
		userMinimal User
	)

	BeforeEach(func() {
		session, err := createSession()
		Expect(err).To(BeNil(), "All tests need a connection to a mongodb.")
		collection = getCollection(session, dbname, collectionname)
		cleanCollection(collection)

		dao = NewUserDao(session, dbname)

		id = bson.NewObjectId()
		name = "myuser"
		worktime, _ = time.ParseDuration("7h45m")
		overtime, _ = time.ParseDuration("1h")
		userWithID = NewPersistedUser(id, name, worktime, overtime)
		userMinimal = NewMinimalUser(name, worktime)
	})

	It("should save a user.", func() {
		Expect(dao.Save(userWithID)).To(Succeed())

		var persistedUser User
		findErr := collection.FindId(id).One(&persistedUser)
		Expect(findErr).To(BeNil())
		Expect(persistedUser).To(Equal(userWithID))
	})

	It("should save a user without overtime", func() {
		Expect(dao.Save(userMinimal)).To(Succeed())

		var persistedUser User
		findErr := collection.Find(bson.M{}).One(&persistedUser)
		Expect(findErr).To(BeNil())
		Expect(persistedUser.EqualsWithoutID(userMinimal)).To(BeTrue())
	})

	It("should find a user by id.", func() {
		Expect(dao.Save(userWithID)).To(Succeed())

		persistedUser, err := dao.FindByID(id)

		Expect(err).To(BeNil())
		Expect(persistedUser).To(Equal(userWithID))
	})

	It("should find a user by name.", func() {
		Expect(dao.Save(userMinimal)).To(Succeed())

		persistedUser, err := dao.FindByName(userMinimal.Name)
		Expect(err).To(BeNil())
		Expect(persistedUser.EqualsWithoutID(userMinimal)).To(BeTrue())
	})

	It("should return error when no user found by name.", func() {
		Expect(dao.Save(userMinimal)).To(Succeed())

		_, err := dao.FindByName("nobody")
		Expect(err.Error()).To(Equal("not found"))
	})

	PIt("should return the password hash from a user.")
	PIt("should save a password hash to a user.")
})
