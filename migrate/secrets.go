package migrate

import (
	"database/sql"
	"fmt"

	"github.com/russross/meddler"
	"github.com/sirupsen/logrus"
)

// MigrateSecrets migrates the secrets V0 database
// to the V1 database.
func MigrateSecrets(source, target *sql.DB) error {
	secretsV0 := []*SecretV0{}

	if err := meddler.QueryAll(source, &secretsV0, secretImportQuery); err != nil {
		return err
	}

	logrus.Infof("migrating %d secrets", len(secretsV0))
	tx, err := target.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	var sequence int64
	for _, secretV0 := range secretsV0 {
		if secretV0.ID > sequence {
			sequence = secretV0.ID
		}

		log := logrus.WithFields(logrus.Fields{
			"repo":   secretV0.RepoID,
			"secret": secretV0.Name,
		})

		log.Debugln("migrate secret")

		secretV1 := &SecretV1{
			ID:     secretV0.ID,
			RepoID: secretV0.RepoID,
			Name:   secretV0.Name,
			Data:   secretV0.Value,
		}

		for _, event := range secretV0.Events {
			if event == "pull_request" {
				secretV1.PullRequest = true
				break
			}
		}

		if err := meddler.Insert(tx, "secrets", secretV1); err != nil {
			log.WithError(err).Errorln("migration failed")
			return err
		}

		log.Debugln("migration complete")
	}

	if meddler.Default == meddler.PostgreSQL {
		_, err = tx.Exec(fmt.Sprintf(updateSecretsSeq, sequence+1))
		if err != nil {
			logrus.WithError(err).Errorln("failed to reset sequence")
			return err
		}
	}

	logrus.Infof("migration complete")
	return tx.Commit()
}

// EncryptSecrets is a helper function that encrypts all database
// secrets after being inserted into the Drone database.
func EncryptSecrets(target *sql.DB, key string) error {
	block, err := parseKey(key)
	if err != nil {
		logrus.WithError(err).Errorln("cannot read encryption key")
		return err
	}

	secretsV1 := []*SecretV1{}

	if err := meddler.QueryAll(target, &secretsV1, secretListQuery); err != nil {
		return err
	}

	logrus.Infof("encrypting %d secrets", len(secretsV1))
	tx, err := target.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	updateStmt := updateSecretStmt
	if meddler.Default == meddler.PostgreSQL {
		updateStmt = updateSecretStmtPostgres
	}

	for _, secretV1 := range secretsV1 {
		ciphertext, err := encrypt(block, secretV1.Data)
		if err != nil {
			logrus.WithError(err).Errorln("encryption failed")
			return err
		}

		secretV1.Data = string(ciphertext)
		if _, err := tx.Exec(updateStmt, secretV1.Data, secretV1.ID); err != nil {
			logrus.WithError(err).Errorln("update failed")
			return err
		}
	}

	logrus.Infof("encryption complete")
	return tx.Commit()
}

const secretListQuery = `
SELECT *
FROM secrets
`

const secretImportQuery = `
SELECT secrets.*
FROM secrets
INNER JOIN repos ON secrets.secret_repo_id = repos.repo_id
WHERE repos.repo_user_id > 0
`

const repoSlugQuery = `
SELECT *
FROM repos
WHERE repo_slug = '%s'
`

const updateSecretsSeq = `
ALTER SEQUENCE secrets_secret_id_seq
RESTART WITH %d
`

const updateSecretStmt = `
UPDATE secrets
SET secret_data = ?
WHERE secret_id = ?
`

const updateSecretStmtPostgres = `
UPDATE secrets
SET secret_data = $1
WHERE secret_id = $2
`
