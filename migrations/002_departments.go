package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func Departments() {
	goose.AddMigration(UpDepartments, DownDepartments)
}

func UpDepartments(tx *sql.Tx) error {
	query := `CREATE TABLE IF NOT EXISTS departments (
						id SERIAL PRIMARY KEY,
						name VARCHAR(200) NOT NULL,
						parent_id INTEGER,
						created_at TIMESTAMP WITH TIME ZONE,
						CHECK (parent_id <> id),
						CONSTRAINT fk_parent_id
							FOREIGN KEY (parent_id)
							REFERENCES departments(id)
							ON DELETE CASCADE
						);`
	_, err := tx.Exec(query)
	return err
}

func DownDepartments(tx *sql.Tx) error {
	query := `DROP TABLE departments;`
	_, err := tx.Exec(query)
	return err
}
