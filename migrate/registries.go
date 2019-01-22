package migrate

import (
	"database/sql"
)

// MigrateRegistries migrates the registry crendeitals
// from the V0 database to the V1 database.
func MigrateRegistries(source, target *sql.DB) error {
	return nil
}
