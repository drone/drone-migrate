package migrate

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/russross/meddler"
	"github.com/sirupsen/logrus"
)

// MigrateRegistries migrates the registry crendeitals
// from the V0 database to the V1 database.
func MigrateRegistries(source, target *sql.DB) error {
	registriesV0 := []*RegistryV0{}
	dockerConfigs := make(map[string]DockerConfig, 0)

	if err := meddler.QueryAll(source, &registriesV0, registryImportQuery); err != nil {
		return err
	}

	logrus.Infof("migrating %d registries", len(registriesV0))
	tx, err := target.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, registryV0 := range registriesV0 {
		log := logrus.WithFields(logrus.Fields{
			"repo": registryV0.RepoFullname,
			"addr": registryV0.Addr,
		})

		log.Debugln("prepare registry")

		if _, ok := dockerConfigs[registryV0.RepoFullname]; !ok {
			dockerConfigs[registryV0.RepoFullname] = DockerConfig{
				AuthConfigs: make(map[string]AuthConfig, 0),
			}
		}

		dockerConfigs[registryV0.RepoFullname].AuthConfigs[registryV0.Addr] = AuthConfig{
			Username: registryV0.Username,
			Password: registryV0.Password,
			Email:    registryV0.Email,
		}

		log.Debugln("prepare complete")
	}

	for repoFullname, dockerConfig := range dockerConfigs {
		log := logrus.WithFields(logrus.Fields{
			"repo": repoFullname,
		})

		log.Debugln("migrate registry")

		result, err := json.Marshal(dockerConfig)

		if err != nil {
			log.WithError(err).Errorln("failed to build docker config")
			continue
		}

		repoV1 := &RepoV1{}

		if err := meddler.QueryRow(target, repoV1, fmt.Sprintf(repoSlugQuery, repoFullname)); err != nil {
			log.WithError(err).Errorln("failed to get registry repo")
			continue
		}

		registryV1 := &RegistryV1{
			RepoID:      repoV1.ID,
			Name:        ".dockerconfigjson",
			Data:        string(result),
			PullRequest: true,
		}

		if err := meddler.Insert(tx, "secrets", registryV1); err != nil {
			log.WithError(err).Errorln("migration failed")
			return err
		}

		log.Debugln("migration complete")
	}

	logrus.Infof("migration complete")
	return tx.Commit()
}

const registryImportQuery = `
SELECT
	repo_full_name,
	registry.*
FROM registry INNER JOIN repos ON (repo_id = registry_repo_id)
WHERE repo_user_id > 0
`
