package dao

import (
	"fmt"
)

type BaseRecord struct {
	Id int
	Sha1 string
	Name string
}

func (r *BaseRecord) Print() {
	fmt.Printf("%d, %s, %s \n", r.Id, r.Sha1, r.Name)
}
