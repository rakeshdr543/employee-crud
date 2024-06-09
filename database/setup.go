package database

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func SetUpDatabase() (*sql.DB, error) {
	// load from env
	connString := "postgresql://root:secret@localhost:5437/employee?sslmode=disable"
	db, err := sql.Open("pgx", connString)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Create the employee table if it doesn't exist
	createEmployeesTable(db)

	return db, nil

}

func createEmployeesTable(db *sql.DB) {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS employee (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		position VARCHAR(100),
		salary Float
	);`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create table: ", err)
		return
	}
}
