package main

import (
	"fmt"
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type BaseDocument struct {
	Id int `json:"id" bson:"_id"`
	Sha1 string
	Name string
}

func (d *BaseDocument) Print() {
	fmt.Printf("%d, %s, %s \n", d.Id, d.Sha1, d.Name)
}

func main() {
	now := time.Now()
	nanos := now.UnixNano()

	query1()

	now2 := time.Now()
	nanos2 := now2.UnixNano()

	fmt.Printf("Took: %d microseconds\n", (nanos2 - nanos) / 1000)
}

func query1() {
	connStr := "dbu:dbp@xmongo:27017/test"
	session, _ := mgo.Dial(connStr)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("file_storage")

	query := []bson.M{
		bson.M{"$match": bson.M{"count": bson.M{"$gt": 0}}},
		bson.M{"$unwind": "$files"},
		bson.M{"$project": bson.M{"_id": 1, "sha1": 1, "name": "$files.name"}},
		bson.M{"$sort": bson.M{"_id": -1, "name": 1}},
		bson.M{"$skip": 1000},
		bson.M{"$limit": 10}}
	var docs []BaseDocument
	c.Pipe(query).All(&docs)

	for _, d := range docs {
		d.Print()
	}
}
