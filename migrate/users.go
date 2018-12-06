package migrate

import (
	"database/sql"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dchest/uniuri"

	"github.com/russross/meddler"
)

// MigrateUsers migrates the user accounts from the V0
// database to the V1 database.
func MigrateUsers(source, target *sql.DB) error {
	usersV0 := []*UserV0{}

	// 1. load all users from the V0 database.
	err := meddler.QueryAll(source, &usersV0, "select * from users")
	if err != nil {
		return err
	}

	// 2. create a database transaction so that we
	// can rollback if the data migration fails.
	tx, err := target.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 3. iterate through the list and convert from
	// the 0.x to the 1.x structure and insert.
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
		err = meddler.Insert(tx, "users", userV1)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
