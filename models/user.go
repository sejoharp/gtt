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

//TODO: Add equals method without ID and test it like the interval one
