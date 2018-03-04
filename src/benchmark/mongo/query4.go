package main

import (
	"fmt"
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MainDocument struct {
	Id int `json:"id" bson:"_id"`
	Sha1 string
	Files []File
}

type File struct {
	Name string
}

func (d *MainDocument) Print() {
	for _, file := range d.Files {
		fmt.Printf("%d, %s, %s \n", d.Id, d.Sha1, file.Name)
	}
}

func main() {
	now := time.Now()
	nanos := now.UnixNano()

	query2()

	now2 := time.Now()
	nanos2 := now2.UnixNano()

	fmt.Printf("Took: %d microseconds\n", (nanos2 - nanos) / 1000)
}

func query2() {
	connStr := "dbu:dbp@xmongo:27017/test"
	session, _ := mgo.Dial(connStr)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("file_storage")

	query := bson.M{"sha1": "806b9a087e6822c1548c606e8e6348b7f08b62ff"}
	projection := bson.M{"_id": 1, "sha1": 1, "files.name": 1}
	var docs []MainDocument
	c.Find(query).Select(projection).All(&docs)

	for _, d := range docs {
		d.Print()
	}
}
