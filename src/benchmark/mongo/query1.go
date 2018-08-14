package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"

	"benchmark/mongo/dao"
)

func Query1() {
	session, c := GetDBC()
	defer session.Close()
	ExecQuery1(c)
}

func ExecQuery1(c *mgo.Collection) {
	query := []bson.M{
		bson.M{"$match": bson.M{"count": bson.M{"$gt": 0}}},
		bson.M{"$unwind": "$files"},
		bson.M{"$project": bson.M{"_id": 1, "sha1": 1, "name": "$files.name"}},
		bson.M{"$sort": bson.M{"_id": -1, "name": 1}},
		bson.M{"$skip": 1000},
		bson.M{"$limit": 10}}
	var docs []dao.BaseDocument
	c.Pipe(query).All(&docs)

	for _, d := range docs {
		d.Print()
	}
}
