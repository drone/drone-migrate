package db

import (
	"database/sql"

	"github.com/drone/drone-migrate/migrate/db/mysql"
	"github.com/drone/drone-migrate/migrate/db/postgres"
	"github.com/drone/drone-migrate/migrate/db/sqlite"
)

// Create creates the 1.0 database.
func Create(db *sql.DB, driver string) error {
	switch driver {
	case "mysql":
		return mysql.Migrate(db)
	case "postgres":
		return postgres.Migrate(db)
	default:
		return sqlite.Migrate(db)
	}
}
