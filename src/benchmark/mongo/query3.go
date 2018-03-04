package main

import (
	"fmt"
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SimpleDocument struct {
	Id int `json:"id" bson:"_id"`
	Sha1 string
	// No file name
}

func (d *SimpleDocument) Print() {
	fmt.Printf("%d, %s \n", d.Id, d.Sha1)
}

func main() {
	now := time.Now()
	nanos := now.UnixNano()

	query3()

	now2 := time.Now()
	nanos2 := now2.UnixNano()

	fmt.Printf("Took: %d microseconds\n", (nanos2 - nanos) / 1000)
}

func query3() {
	connStr := "dbu:dbp@xmongo:27017/test"
	session, _ := mgo.Dial(connStr)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("file_storage")

	query := bson.M{"_id": bson.M{"$in": [2]int{171, 352}}}
	projection := bson.M{"_id": 1, "sha1": 1}
	var docs []SimpleDocument
	c.Find(query).Select(projection).All(&docs)

	for _, d := range docs {
		d.Print()
	}
}
