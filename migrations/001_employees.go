package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func Employees() {
	goose.AddMigration(UpEmployees, DownEmployees)
}

func UpEmployees(tx *sql.Tx) error {
	query := `CREATE TABLE IF NOT EXISTS employees (
						id 							SERIAL PRIMARY KEY,
						department_id 	bigint,
						full_name				varchar(200) NOT NULL,
						position 				varchar(200) NOT NULL,
						hired_at				date,
						created_at			timestamp with time zone
						);`
	_, err := tx.Exec(query)
	return err
}

func DownEmployees(tx *sql.Tx) error {
	query := `DROP TABLE employees;`
	_, err := tx.Exec(query)
	return err
}
