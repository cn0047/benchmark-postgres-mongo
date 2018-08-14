package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"

	"benchmark/mongo/dao"
)

func Query3() {
	session, c := GetDBC()
	defer session.Close()
	ExecQuery3(c)
}

func ExecQuery3(c *mgo.Collection) {
	query := bson.M{"_id": bson.M{"$in": [...]int{171, 352}}}
	projection := bson.M{"_id": 1, "sha1": 1}
	var docs []dao.SimpleDocument
	c.Find(query).Select(projection).All(&docs)

	for _, d := range docs {
		d.Print()
	}
}
