package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dchest/uniuri"
	"github.com/drone/drone-go/drone"
	"github.com/drone/go-scm/scm"
	"github.com/hashicorp/go-multierror"
	"github.com/russross/meddler"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// MigrateRepos migrates the repositories from the V0
// database to the V1 database.
func MigrateRepos(source, target *sql.DB) error {
	reposV0 := []*RepoV0{}

	if err := meddler.QueryAll(source, &reposV0, repoImportQuery); err != nil {
		return err
	}

	logrus.Infof("migrating %d repositories", len(reposV0))

	tx, err := target.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	var sequence int64
	for _, repoV0 := range reposV0 {
		if repoV0.ID > sequence {
			sequence = repoV0.ID
		}

		log := logrus.WithFields(logrus.Fields{
			"id":   repoV0.ID,
			"repo": repoV0.FullName,
		})

		log.Debugln("migrate repository")

		repoV1 := &RepoV1{
			ID:         repoV0.ID,
			UserID:     repoV0.UserID,
			Namespace:  repoV0.Owner,
			Name:       repoV0.Name,
			Slug:       repoV0.FullName,
			SCM:        "git",
			HTTPURL:    repoV0.Clone,
			SSHURL:     "",
			Link:       repoV0.Link,
			Branch:     repoV0.Branch,
			Private:    repoV0.IsPrivate,
			Visibility: repoV0.Visibility,
			Active:     repoV0.IsActive,
			Config:     repoV0.Config,
			Trusted:    repoV0.IsTrusted,
			Protected:  repoV0.IsGated,
			Timeout:    repoV0.Timeout,
			Counter:    int64(repoV0.Counter),
			Synced:     time.Now().Unix(),
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Version:    1,
			Signer:     uniuri.NewLen(32),
			Secret:     uniuri.NewLen(32),

			// We use a temporary repository identifier here.
			// We need to do a per-repository lookup to get the
			// actual identifier from the source code management
			// system.
			UID: fmt.Sprintf("temp_%d", repoV0.ID),
		}

		if err := meddler.Insert(tx, "repos", repoV1); err != nil {
			log.WithError(err).Errorln("migration failed")
			return err
		}

		log.Debugln("migration complete")
	}

	if meddler.Default == meddler.PostgreSQL {
		_, err = tx.Exec(fmt.Sprintf(updateRepoSeq, sequence+1))
		if err != nil {
			logrus.WithError(err).Errorln("failed to reset sequence")
			return err
		}
	}

	logrus.Infoln("migration complete")
	return tx.Commit()
}

// UpdateRepoIdentifiers updates the repository identifiers
// from temporary values (assigned during migration) to the
// value fetched from the source code management system.
func UpdateRepoIdentifiers(db *sql.DB, client *scm.Client) error {
	repos := []*RepoV1{}
	var result error

	if err := meddler.QueryAll(db, &repos, repoTempQuery); err != nil {
		return err
	}

	logrus.Infoln("updating repository metadata")

	for _, repo := range repos {
		log := logrus.WithFields(logrus.Fields{
			"repo": repo.Slug,
		})

		user := &UserV1{}

		if err := meddler.QueryRow(db, user, fmt.Sprintf(userIdentifierQuery, repo.UserID)); err != nil {
			log.WithError(err).Errorf("failed to get repository owner")
			multierror.Append(result, err)
			continue
		}

		log = log.WithField("owner", user.Login)

		tok := &scm.Token{
			Token:   user.Token,
			Refresh: user.Refresh,
		}
		if user.Expiry > 0 {
			tok.Expires = time.Unix(user.Expiry, 0)
		}
		ctx := scm.WithContext(context.Background(), tok)

		remoteRepo, _, err := client.Repositories.Find(ctx, scm.Join(repo.Namespace, repo.Name))

		if err != nil {
			log.WithError(err).Errorf("failed to get remote repository")
			multierror.Append(result, err)
			continue
		}

		if _, err := db.Exec(fmt.Sprintf(repoUpdateQuery, remoteRepo.ID, repo.ID)); err != nil {
			log.WithError(err).Errorf("failed to update metadata")
			multierror.Append(result, err)
		}

		log.Debugln("updated metadata")
	}

	logrus.Infoln("repository metadata update complete")
	return result
}

// ActivateRepositories re-activates the repositories.
// This will create new webhooks and populate any empty
// values (security keys, etc).
func ActivateRepositories(db *sql.DB, client drone.Client) error {
	repos := []*RepoV1{}
	var result error

	if err := meddler.QueryAll(db, &repos, repoActivateQuery); err != nil {
		return err
	}

	logrus.Infoln("begin repository activation")

	for _, repo := range repos {
		log := logrus.WithFields(logrus.Fields{
			"repo": repo.Slug,
		})

		log.Debugln("activating repository")

		user := &UserV1{}

		if err := meddler.QueryRow(db, user, fmt.Sprintf(userIdentifierQuery, repo.UserID)); err != nil {
			log.WithError(err).Errorf("failed to get repository owner")
			multierror.Append(result, err)
			continue
		}

		log = log.WithField("owner", user.Login)

		config := new(oauth2.Config)

		client.SetClient(config.Client(
			oauth2.NoContext,
			&oauth2.Token{
				AccessToken:  user.Token,
				RefreshToken: user.Refresh,
			},
		))

		if _, err := client.RepoPost(repo.Namespace, repo.Name); err != nil {
			log.WithError(err).Errorf("activation failed")
			multierror.Append(result, err)
			continue
		}

		log.Debugln("successfully activated")
	}

	logrus.Infoln("repository activation complete")
	return result
}

const repoImportQuery = `
SELECT *
FROM repos
WHERE repo_user_id > 0
`

const repoTempQuery = `
SELECT *
FROM repos
WHERE repo_uid LIKE 'temp_%'
`

const userIdentifierQuery = `
SELECT *
FROM users
WHERE user_id = %d
`

const repoUpdateQuery = `
UPDATE repos
SET repo_uid = '%s'
WHERE repo_id = %d
`

const repoActivateQuery = `
SELECT *
FROM repos
`

const updateRepoSeq = `
ALTER SEQUENCE repos_repo_id_seq
RESTART WITH %d
`
