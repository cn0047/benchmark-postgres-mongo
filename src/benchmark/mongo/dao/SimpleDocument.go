package dao

import (
	"fmt"
)

type SimpleDocument struct {
	Id int `json:"id" bson:"_id"`
	Sha1 string
	// No file name
}

func (d *SimpleDocument) Print() {
	fmt.Printf("%d, %s \n", d.Id, d.Sha1)
}
