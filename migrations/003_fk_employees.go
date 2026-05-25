package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func EmployeesFK() {
	goose.AddMigration(UpEmployeesFK, DownEmployeesFK)
}

func UpEmployeesFK(tx *sql.Tx) error {
	query := `ALTER TABLE employees
						ADD CONSTRAINT fk_employees_departments
						FOREIGN KEY (department_id)
						REFERENCES departments(id)
						ON DELETE CASCADE;
						`
	_, err := tx.Exec(query)
	return err
}

func DownEmployeesFK(tx *sql.Tx) error {
	query := `ALTER TABLE employees DROP CONSTRAINT fk_employees_departments;`
	_, err := tx.Exec(query)
	return err
}
