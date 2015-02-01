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
	return Interval{UserID: userID, Start: cleanNanoSeconds(start), Stop: cleanNanoSeconds(stop)}
}

func NewIntervalWithStart(userID bson.ObjectId, start time.Time) Interval {
	return Interval{UserID: userID, Start: cleanNanoSeconds(start)}
}

func NewPersistedInterval(id bson.ObjectId, userID bson.ObjectId, start time.Time, stop time.Time) Interval {
	return Interval{id, userID, cleanNanoSeconds(start), cleanNanoSeconds(stop)}
}

func NewPersistedIntervalWithStart(id bson.ObjectId, userID bson.ObjectId, start time.Time) Interval {
	return Interval{ID: id, UserID: userID, Start: cleanNanoSeconds(start)}
}

func cleanNanoSeconds(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(),
		date.Hour(), date.Minute(), date.Second(), 0, date.Location())
}

func (interval Interval) EqualsWithOutID(that Interval) bool {
	if interval.UserID == that.UserID && interval.Start == that.Start && interval.Stop == that.Stop {
		return true
	}
	return false
}
