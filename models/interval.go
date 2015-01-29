package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//Interval contains the Duration of a work from a user
type Interval struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	UserID bson.ObjectId
	Start  time.Time
	Stop   time.Time `bson:",omitempty"`
}

func NewInterval(userID bson.ObjectId, start time.Time, stop time.Time) Interval {
	return Interval{UserID: userID, Start: start, Stop: stop}
}

func NewIntervalWithStart(userID bson.ObjectId, start time.Time) Interval {
	return Interval{UserID: userID, Start: start}
}

func NewPersistedInterval(id bson.ObjectId, userID bson.ObjectId, start time.Time, stop time.Time) Interval {
	return Interval{id, userID, start, stop}
}

func NewPersistedIntervalWithStart(userID bson.ObjectId, start time.Time) Interval {
	return Interval{UserID: userID, Start: start}
}
