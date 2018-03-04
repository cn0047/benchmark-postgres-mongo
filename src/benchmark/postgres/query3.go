package main

import (
	"fmt"
	"time"
	"database/sql"
	_ "github.com/lib/pq"
)

type SimpleRecord struct {
	Id int
	Sha1 string
	// No file name
}

func (r *SimpleRecord) Print() {
	fmt.Printf("%d, %s \n", r.Id, r.Sha1)
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
	connStr := "host=xpostgres port=5432 user=dbu password=dbp dbname=test sslmode=disable"
	db, _ := sql.Open("postgres", connStr)
	defer db.Close()

	query := `
		SELECT s.id, s.sha1
		FROM storage s
		WHERE s.id IN (171, 352)
	`
	rows, _ := db.Query(query)
	defer rows.Close()

	r := &SimpleRecord{}
	for rows.Next() {
		rows.Scan(&r.Id, &r.Sha1)
		r.Print()
	}
}
