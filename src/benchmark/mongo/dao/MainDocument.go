package dao

import (
	"fmt"
)

type MainDocument struct {
	Id int `json:"id" bson:"_id"`
	Sha1 string
	Files []FileDocument
}

func (d *MainDocument) Print() {
	for _, file := range d.Files {
		fmt.Printf("%d, %s, %s \n", d.Id, d.Sha1, file.Name)
	}
}
