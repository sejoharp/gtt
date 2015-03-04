package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string
	Worktime time.Duration
	Overtime time.Duration
}

type UserWithPassword struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string
	Worktime time.Duration
	Overtime time.Duration
	Password string
}

func NewMinimalUserWithPassword(name string, worktime time.Duration, password string) UserWithPassword {
	return UserWithPassword{Name: name, Worktime: worktime, Overtime: 0, Password: password}
}

func NewMinimalUser(name string, worktime time.Duration) User {
	return User{Name: name, Worktime: worktime, Overtime: 0}
}

func NewPersistedMinimalUser(id bson.ObjectId, name string, worktime time.Duration) User {
	return User{ID: id, Name: name, Worktime: worktime, Overtime: 0}
}

func NewUser(name string, worktime time.Duration, ovetime time.Duration) User {
	return User{Name: name, Worktime: worktime, Overtime: ovetime}
}

func NewPersistedUser(id bson.ObjectId, name string, worktime time.Duration, ovetime time.Duration) User {
	return User{ID: id, Name: name, Worktime: worktime, Overtime: ovetime}
}

func (user User) EqualsWithoutID(that User) bool {
	if user.Name == that.Name && //
		user.Overtime == that.Overtime && //
		user.Worktime == that.Worktime {
		return true
	}
	return false
}

func (user UserWithPassword) EqualsWithoutID(that UserWithPassword) bool {
	if user.Name == that.Name && //
		user.Overtime == that.Overtime && //
		user.Worktime == that.Worktime && //
		user.Password == that.Password {
		return true
	}
	return false
}
