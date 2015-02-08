package models

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("UserDao", func() {
	const dbName = "timetracker"
	const collectionName = "users"
	const passwordHash = "3!aYBlA994"

	var (
		collection             *mgo.Collection
		dao                    *UserDao
		id                     bson.ObjectId
		name                   string
		worktime               time.Duration
		overtime               time.Duration
		userPersistedWithID    User
		userPersistedWithoutID User
		userWithoutID          User
	)

	BeforeEach(func() {
		session, err := createSession()
		Expect(err).To(BeNil(), "All tests need a connection to a mongodb.")
		collection = getCollection(session, dbName, collectionName)
		cleanCollection(collection)

		dao = NewUserDao(session, dbName)

		id = bson.NewObjectId()
		name = "myuser"
		worktime, _ = time.ParseDuration("7h45m")
		overtime, _ = time.ParseDuration("1h")
		userPersistedWithID = NewPersistedUser(id, name, worktime, overtime)
		userPersistedWithoutID = NewPersistedMinimalUser(id, name, worktime)
		userWithoutID = NewMinimalUser(name, worktime)
	})

	It("should save a user.", func() {
		Expect(dao.Save(userPersistedWithID)).To(Succeed())

		var persistedUser User
		findErr := collection.FindId(id).One(&persistedUser)
		Expect(findErr).To(BeNil())
		Expect(persistedUser).To(Equal(userPersistedWithID))
	})

	It("should save a user without overtime", func() {
		Expect(dao.Save(userWithoutID)).To(Succeed())

		var persistedUser User
		findErr := collection.Find(bson.M{}).One(&persistedUser)
		Expect(findErr).To(BeNil())
		Expect(persistedUser.EqualsWithoutID(userWithoutID)).To(BeTrue())
	})

	It("should find a user by id.", func() {
		Expect(dao.Save(userPersistedWithoutID)).To(Succeed())

		persistedUser, err := dao.FindByID(id)

		Expect(err).To(BeNil())
		Expect(persistedUser).To(Equal(userPersistedWithoutID))
	})

	It("should find a user by name.", func() {
		Expect(dao.Save(userWithoutID)).To(Succeed())

		persistedUser, err := dao.FindByName(userWithoutID.Name)
		Expect(err).To(BeNil())
		Expect(persistedUser.EqualsWithoutID(userWithoutID)).To(BeTrue())
	})

	It("should return error when no user found by name.", func() {
		Expect(dao.Save(userWithoutID)).To(Succeed())

		_, err := dao.FindByName("nobody")
		Expect(err.Error()).To(Equal("not found"))
	})

	It("should save a password hash to a user.", func() {
		Expect(dao.Save(userPersistedWithID)).To(Succeed())

		Expect(dao.AddPassword(userPersistedWithID.ID, passwordHash)).To(Succeed())

		var result bson.M
		err := collection.FindId(userPersistedWithID.ID).Select(bson.M{"password": 1}).One(&result)

		Expect(err).To(BeNil())
		Expect(result["password"]).To(Equal(passwordHash))
	})

	It("should return the password hash from a user.", func() {
		Expect(dao.Save(userPersistedWithID)).To(Succeed())
		Expect(dao.AddPassword(userPersistedWithID.ID, passwordHash)).To(Succeed())

		hash, err := dao.GetPassword(userPersistedWithID.ID)

		Expect(err).To(BeNil())
		Expect(hash).To(Equal(passwordHash))
	})
})
