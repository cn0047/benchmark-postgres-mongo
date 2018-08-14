package postgres

import (
	_ "github.com/lib/pq"
	"database/sql"

	"benchmark/postgres/dao"
)

func Query3() {
	db := GetDBC()
	defer db.Close()
	ExecQuery3(db)
}

func ExecQuery3(db *sql.DB) {
	query := `
		SELECT s.id, s.sha1
		FROM storage s
		WHERE s.id IN (171, 352)
	`
	rows, _ := db.Query(query)
	defer rows.Close()

	r := &dao.SimpleRecord{}
	for rows.Next() {
		rows.Scan(&r.Id, &r.Sha1)
		r.Print()
	}
}
