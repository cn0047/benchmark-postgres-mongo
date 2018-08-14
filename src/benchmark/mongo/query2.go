package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"

	"benchmark/mongo/dao"
)

func Query2() {
	session, c := GetDBC()
	defer session.Close()
	ExecQuery2(c)
}

func ExecQuery2(c *mgo.Collection) {
	query := bson.M{"sha1": "806b9a087e6822c1548c606e8e6348b7f08b62ff"}
	projection := bson.M{"_id": 1, "sha1": 1, "files": bson.M{"$elemMatch": bson.M{"name": "Sit.avi"}}}
	var docs []dao.MainDocument
	c.Find(query).Select(projection).All(&docs)

	for _, d := range docs {
		d.Print()
	}
}
