package mongo

import (
	"gopkg.in/mgo.v2"
)

func GetDBC() (*mgo.Session, *mgo.Collection) {
	connStr := "dbu:dbp@xmongo:27017/test"
	session, _ := mgo.Dial(connStr)
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("file_storage")

	return session, c
}