package dao

import (
	"fmt"
)

type BaseDocument struct {
	Id int `json:"id" bson:"_id"`
	Sha1 string
	Name string
}

func (d *BaseDocument) Print() {
	fmt.Printf("%d, %s, %s \n", d.Id, d.Sha1, d.Name)
}
