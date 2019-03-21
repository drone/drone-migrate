package migrate

import (
	"database/sql"
	"time"

	"github.com/dchest/uniuri"
	"github.com/sirupsen/logrus"

	"github.com/russross/meddler"
)

// MigrateUsers migrates the user accounts from the V0
// database to the V1 database.
func MigrateUsers(source, target *sql.DB) error {
	usersV0 := []*UserV0{}

	if err := meddler.QueryAll(source, &usersV0, userImportQuery); err != nil {
		return err
	}

	logrus.Infof("migrating %d users", len(usersV0))

	tx, err := target.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, userV0 := range usersV0 {
		log := logrus.WithFields(logrus.Fields{
			"id":    userV0.ID,
			"login": userV0.Login,
		})

		log.Debugln("migrate user")

		userV1 := &UserV1{
			ID:        userV0.ID,
			Login:     userV0.Login,
			Email:     userV0.Email,
			Machine:   false,
			Admin:     false,
			Active:    true,
			Avatar:    userV0.Avatar,
			Syncing:   false,
			Synced:    0,
			Created:   time.Now().Unix(),
			Updated:   time.Now().Unix(),
			LastLogin: 0,
			Token:     userV0.Token,
			Refresh:   userV0.Secret,
			Expiry:    userV0.Expiry,
			Hash:      uniuri.NewLen(32),
		}

		if err := meddler.Insert(tx, "users", userV1); err != nil {
			log.WithError(err).Errorln("migration failed")
			return err
		}

		log.Debugln("migration complete")
	}

	logrus.Infoln("migration complete")
	return tx.Commit()
}

const userImportQuery = `
SELECT
	*
FROM
	users
`
