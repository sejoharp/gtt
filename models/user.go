package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Worktime time.Duration
	Overtime time.Duration
}

func NewMinimalUser(worktime time.Duration) User {
	return User{Worktime: worktime, Overtime: 0}
}

func NewPersistedMinimalUser(id bson.ObjectId, worktime time.Duration) User {
	return User{ID: id, Worktime: worktime, Overtime: 0}
}

func NewUser(worktime time.Duration, ovetime time.Duration) User {
	return User{Worktime: worktime, Overtime: ovetime}
}

func NewPersistedUser(id bson.ObjectId, worktime time.Duration, ovetime time.Duration) User {
	return User{ID: id, Worktime: worktime, Overtime: ovetime}
}
