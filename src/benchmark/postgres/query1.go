package main

import (
	"fmt"
	"time"
	"database/sql"
	_ "github.com/lib/pq"
)

type BaseRecord struct {
	Id int
	Sha1 string
	Name string
}

func (r *BaseRecord) Print() {
	fmt.Printf("%d, %s, %s \n", r.Id, r.Sha1, r.Name)
}

func main() {
	now := time.Now()
	nanos := now.UnixNano()

	query1()

	now2 := time.Now()
	nanos2 := now2.UnixNano()

	fmt.Printf("Took: %d microseconds\n", (nanos2 - nanos) / 1000)
}

func query1() {
	connStr := "host=xpostgres port=5432 user=dbu password=dbp dbname=test sslmode=disable"
	db, _ := sql.Open("postgres", connStr)
	defer db.Close()

	query := `
		SELECT s.id, s.sha1, f.name
		FROM storage s
		JOIN file f ON s.id = f.storage_id
		WHERE s.count > 0
		ORDER by id DESC, name ASC
		OFFSET 1000 LIMIT 10
	`
	rows, _ := db.Query(query)
	defer rows.Close()

	r := &BaseRecord{}
	for rows.Next() {
		rows.Scan(&r.Id, &r.Sha1, &r.Name)
		r.Print()
	}
}
