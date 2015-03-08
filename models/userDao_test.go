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
		collection       *mgo.Collection
		dao              UserDao
		id               bson.ObjectId
		name             string
		worktime         time.Duration
		overtime         time.Duration
		oldUser          User
		oldUserWithoutID User
		newUserWithoutID User
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
		oldUser = NewPersistedUser(id, name, worktime, overtime)
		oldUserWithoutID = NewPersistedMinimalUser(id, name, worktime)
		newUserWithoutID = NewMinimalUser(name, worktime)
	})

	It("should save a user.", func() {
		Expect(dao.Save(oldUser)).To(Succeed())

		var persistedUser User
		findErr := collection.FindId(id).One(&persistedUser)
		Expect(findErr).To(BeNil())
		Expect(persistedUser).To(Equal(oldUser))
	})

	It("should save a user without overtime", func() {
		Expect(dao.Save(newUserWithoutID)).To(Succeed())

		var persistedUser User
		findErr := collection.Find(bson.M{}).One(&persistedUser)
		Expect(findErr).To(BeNil())
		Expect(persistedUser.EqualsWithoutID(newUserWithoutID)).To(BeTrue())
	})

	It("should find a user by id.", func() {
		Expect(dao.Save(oldUserWithoutID)).To(Succeed())

		persistedUser, err := dao.FindByID(id)

		Expect(err).To(BeNil())
		Expect(persistedUser).To(Equal(oldUserWithoutID))
	})

	It("should find a user by name.", func() {
		Expect(dao.Save(newUserWithoutID)).To(Succeed())

		persistedUser, err := dao.FindByName(newUserWithoutID.Name)
		Expect(err).To(BeNil())
		Expect(persistedUser.EqualsWithoutID(newUserWithoutID)).To(BeTrue())
	})

	It("should return error when no user found by name.", func() {
		Expect(dao.Save(newUserWithoutID)).To(Succeed())

		_, err := dao.FindByName("nobody")
		Expect(err.Error()).To(Equal("not found"))
	})

	It("should save a password hash to a user.", func() {
		Expect(dao.Save(oldUser)).To(Succeed())

		Expect(dao.AddPassword(oldUser.ID, passwordHash)).To(Succeed())

		var result bson.M
		err := collection.FindId(oldUser.ID).Select(bson.M{"password": 1}).One(&result)
		Expect(err).To(BeNil())
		Expect(result["password"]).To(Equal(passwordHash))
	})

	It("should add a password hash to a user.", func() {
		Expect(dao.Save(oldUser)).To(Succeed())

		Expect(dao.AddPasswordByUser(oldUser.Name, passwordHash)).To(Succeed())

		var result bson.M
		err := collection.FindId(oldUser.ID).Select(bson.M{"password": 1}).One(&result)
		Expect(err).To(BeNil())
		Expect(result["password"]).To(Equal(passwordHash))
	})

	It("should return the password hash from a user.", func() {
		Expect(dao.Save(oldUser)).To(Succeed())
		Expect(dao.AddPassword(oldUser.ID, passwordHash)).To(Succeed())

		hash, err := dao.GetPassword(oldUser.ID)

		Expect(err).To(BeNil())
		Expect(hash).To(Equal(passwordHash))
	})

	It("should save a user with password.", func() {
		user := NewMinimalUserWithPassword(name, worktime, []byte("password"))

		Expect(dao.SaveWithPassword(user)).To(Succeed())

		var persistedUser UserWithPassword
		err := collection.Find(bson.M{"name": user.Name}).One(&persistedUser)
		Expect(err).To(BeNil())
		Expect(persistedUser.EqualsWithoutID(user)).To(BeTrue())
	})

	It("should change the overtime of a user.", func() {
		dao.Save(newUserWithoutID)
		persistedUser, _ := dao.FindByName(newUserWithoutID.Name)

		persistedUser.Name = "newName"
		Expect(dao.Update(persistedUser)).To(Succeed())

		changedUser, updateErr := dao.FindByID(persistedUser.ID)
		Expect(updateErr).To(BeNil())

		Expect(changedUser.EqualsWithoutID(persistedUser)).To(BeTrue())
	})

	It("should return the password hash from a user.", func() {
		user := NewMinimalUserWithPassword(name, worktime, []byte("password"))
		Expect(dao.SaveWithPassword(user)).To(Succeed())

		hash, err := dao.GetPasswordByUser((user.Name))

		Expect(err).To(BeNil())
		Expect(hash).To(Equal(user.Password))
	})

	It("should return an error when user does not exist.", func() {
		user := NewMinimalUserWithPassword(name, worktime, []byte("password"))
		Expect(dao.SaveWithPassword(user)).To(Succeed())

		_, err := dao.GetPasswordByUser(("myname"))

		Expect(err).To(HaveOccurred())
	})
})
