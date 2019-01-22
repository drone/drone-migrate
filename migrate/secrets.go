package migrate

import (
	"database/sql"
)

// MigrateSecrets migrates the secrets V0 database
// to the V1 database.
func MigrateSecrets(source, target *sql.DB) error {
	return nil
}
