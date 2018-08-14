package dao

import (
	"fmt"
)

type SimpleRecord struct {
	Id int
	Sha1 string
	// No file name
}

func (r *SimpleRecord) Print() {
	fmt.Printf("%d, %s \n", r.Id, r.Sha1)
}

