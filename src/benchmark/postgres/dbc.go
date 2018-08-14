package postgres

import "database/sql"

func GetDBC() *sql.DB {
	connStr := "host=xpostgres port=5432 user=dbu password=dbp dbname=test sslmode=disable"
	db, _ := sql.Open("postgres", connStr)

	return db
}
