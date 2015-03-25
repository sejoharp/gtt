package controller

import (
	"github.com/zenazn/goji/web"
	"gopkg.in/mgo.v2/bson"
)

func getUserID(c web.C) bson.ObjectId {
	return c.Env["userID"].(bson.ObjectId)
}
