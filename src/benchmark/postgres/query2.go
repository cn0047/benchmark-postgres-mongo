package main

import (
	"fmt"
	"time"
	"database/sql"
	_ "github.com/lib/pq"
)

type MainRecord struct {
	Id int
	Sha1 string
	Name string
}

func (r *MainRecord) Print() {
	fmt.Printf("%d, %s, %s \n", r.Id, r.Sha1, r.Name)
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
	connStr := "host=xpostgres port=5432 user=dbu password=dbp dbname=test sslmode=disable"
	db, _ := sql.Open("postgres", connStr)
	defer db.Close()

	query := `
		SELECT s.id, s.sha1, f.name
		FROM storage s
		JOIN file f ON s.id = f.storage_id
		WHERE s.sha1 = '806b9a087e6822c1548c606e8e6348b7f08b62ff' AND f.name = 'Sit.avi'
	`
	rows, _ := db.Query(query)
	defer rows.Close()

	r := &MainRecord{}
	for rows.Next() {
		rows.Scan(&r.Id, &r.Sha1, &r.Name)
		r.Print()
	}
}
