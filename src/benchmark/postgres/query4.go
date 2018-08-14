package postgres

import (
	_ "github.com/lib/pq"
	"database/sql"

	"benchmark/postgres/dao"
)

func Query4() {
	db := GetDBC()
	defer db.Close()
	ExecQuery4(db)
}

func ExecQuery4(db *sql.DB) {
	query := `
		SELECT s.id, s.sha1, f.name
		FROM storage s
		JOIN file f ON s.id = f.storage_id
		WHERE s.sha1 = '806b9a087e6822c1548c606e8e6348b7f08b62ff'
	`
	rows, _ := db.Query(query)
	defer rows.Close()

	r := &dao.MainRecord{}
	for rows.Next() {
		rows.Scan(&r.Id, &r.Sha1, &r.Name)
		r.Print()
	}
}
