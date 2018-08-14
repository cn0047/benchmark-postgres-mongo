package dao

import (
	"fmt"
)

type MainRecord struct {
	Id int
	Sha1 string
	Name string
}

func (r *MainRecord) Print() {
	fmt.Printf("%d, %s, %s \n", r.Id, r.Sha1, r.Name)
}
