package migrate

import (
	"database/sql"
	"encoding/json"
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

	for _, secretV0 := range secretsV0 {
		log := logrus.WithFields(logrus.Fields{
			"repo":   secretV0.RepoFullname,
			"secret": secretV0.Name,
		})

		log.Debugln("migrate secret")

		repoV1 := &RepoV1{}

		if err := meddler.QueryRow(target, repoV1, fmt.Sprintf(repoSlugQuery, secretV0.RepoFullname)); err != nil {
			log.WithError(err).Errorln("failed to get secret repo")
			continue
		}

		pullRequest := false

		if secretV0.Events != "" {
			events := make([]string, 0)
			json.Unmarshal([]byte(secretV0.Events), &events)

			for _, event := range events {
				if event == "pull_request" {
					pullRequest = true
					break
				}
			}
		}

		secretV1 := &SecretV1{
			ID:          secretV0.ID,
			RepoID:      repoV1.ID,
			Name:        secretV0.Name,
			Data:        secretV0.Value,
			PullRequest: pullRequest,
		}

		if err := meddler.Insert(tx, "secrets", secretV1); err != nil {
			log.WithError(err).Errorln("migration failed")
			return err
		}

		log.Debugln("migration complete")
	}

	logrus.Infof("migration complete")
	return tx.Commit()
}

const secretImportQuery = `
SELECT
	repo_full_name,
	secrets.*
FROM
	secrets
	INNER JOIN repos ON (repo_id = secret_repo_id)
WHERE
	secret_repo_id > 0
`

const repoSlugQuery = `
SELECT
	*
FROM
	repos
WHERE
	repo_slug = '%s'
`
