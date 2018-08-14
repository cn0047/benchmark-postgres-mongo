package postgres

import (
	_ "github.com/lib/pq"
	"database/sql"

	"benchmark/postgres/dao"
)

func Query1() {
	db := GetDBC()
	defer db.Close()
	ExecQuery1(db)
}

func ExecQuery1(db *sql.DB) {
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

	r := &dao.BaseRecord{}
	for rows.Next() {
		rows.Scan(&r.Id, &r.Sha1, &r.Name)
		r.Print()
	}
}
