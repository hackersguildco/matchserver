package models

import (
	"os"

	"github.com/cheersappio/matchserver/utils"

	"gopkg.in/mgo.v2"
)

var (
	StrokesCollection *mgo.Collection
	Session           *mgo.Session
)

func init() {
	utils.Log.Infof("init DB: %s", os.Getenv("MONGO_URI"))
	session, err := mgo.Dial(os.Getenv("MONGO_URI"))
	if err != nil {
		panic(err)
	}
	Session = session
	Session.SetMode(mgo.Monotonic, true)
	StrokesCollection = Session.DB(os.Getenv("MONGO_DB")).C("strokes")
	index := mgo.Index{
		Key:  []string{"$2dsphere:location"},
		Bits: 26,
	}
	if err := StrokesCollection.EnsureIndex(index); err != nil {
		panic(err)
	}
}
