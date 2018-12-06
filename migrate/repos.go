package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dchest/uniuri"
	"github.com/drone/drone-go/drone"
	"github.com/drone/go-scm/scm"
	"github.com/hashicorp/go-multierror"
	"github.com/russross/meddler"
	"golang.org/x/oauth2"
)

var noContext = context.Background()

// MigrateRepos migrates the repositories from the V0
// database to the V1 database.
func MigrateRepos(source, target *sql.DB) error {
	reposV0 := []*RepoV0{}

	// 1. load all repos from the V0 database.
	err := meddler.QueryAll(source, &reposV0, "select * from repos where repo_user_id > 0")
	if err != nil {
		return err
	}

	logrus.Infof("migrating %d repositories", len(reposV0))

	// 2. create a database transaction so that we
	// can rollback if the data migration fails.
	tx, err := target.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 3. iterate through the list and convert from
	// the 0.x to the 1.x structure and insert.
	for _, repoV0 := range reposV0 {
		log := logrus.WithField("repository", repoV0.FullName)
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
		}

		// We use a temporary repository identifier here.
		// We need to do a per-repository lookup to get the
		// actual identifier from the source code management
		// system.
		repoV1.UID = fmt.Sprintf("temp_%d", repoV0.ID)

		err = meddler.Insert(tx, "repos", repoV1)
		if err != nil {
			log.WithError(err).Errorln("migration failed")
			return err
		}

		log.Debugln("migration complete")
	}

	logrus.Infof("migration complete")
	return tx.Commit()
}

// UpdateRepoIdentifiers updates the repository identifiers
// from temporary values (assigned during migration) to the
// value fetched from the source code management system.
func UpdateRepoIdentifiers(db *sql.DB, client *scm.Client) error {
	repos := []*RepoV1{}

	// 1. load all repos from the V1 database.
	err := meddler.QueryAll(db, &repos, "select * from repos where repo_uid LIKE 'temp_%'")
	if err != nil {
		return err
	}

	logrus.Infof("updating repository metadata")

	var result error
	for _, repo := range repos {
		log := logrus.WithField("repository", repo.Slug)
		log.Debugln("update metadata")

		// 2.a fetch the repository owner
		user := &UserV1{}
		err = meddler.QueryRow(db, user, fmt.Sprintf("SELECT * FROM users WHERE user_id = %d", repo.UserID))
		if err != nil {
			log.WithError(err).Errorf("failed to get repository owner")
			multierror.Append(result, err)
			continue
		}

		log = logrus.WithField("owner", user.Login)

		// 2.b fetch the remote repository by name.
		remoteRepo, _, err := client.Repositories.Find(noContext, scm.Join(repo.Namespace, repo.Name))
		if err != nil {
			log.WithError(err).Errorf("failed to get remote repository")
			multierror.Append(result, err)
			continue
		}

		// 2.c update the temporary id for the remote
		// repository with the value from the remote
		// system.
		repo.UID = remoteRepo.ID
		repo.SSHURL = remoteRepo.CloneSSH
		err = meddler.Update(db, "repos", repo)
		if err != nil {
			log.WithError(err).Errorf("failed to update metadata")
			multierror.Append(result, err)
		}

		log.Debugln("successfully updated metadata")
	}

	logrus.Infoln("repository metadata update complete")
	return result
}

// ActivateRepositories re-activates the repositories.
// This will create new webhooks and populate any empty
// values (security keys, etc).
func ActivateRepositories(db *sql.DB, client drone.Client) error {
	logrus.Infoln("begin repository activation")

	repos := []*RepoV1{}

	err := meddler.QueryAll(db, &repos, "select * from repos")
	if err != nil {
		return err
	}

	var result error
	for _, repo := range repos {
		log := logrus.WithField("repository", repo.Slug)
		log.Debugln("activating repository")

		// 2.a fetch the repository owner
		user := &UserV1{}
		err = meddler.QueryRow(db, user, fmt.Sprintf("SELECT * FROM users WHERE user_id = %d", repo.UserID))
		if err != nil {
			log.WithError(err).Errorf("failed to get repository owner")
			multierror.Append(result, err)
			continue
		}

		log = logrus.WithField("owner", user.Login)

		// 2.b configure the drone client to use the
		// authorization token for the previously fetched
		// user.
		config := new(oauth2.Config)
		auther := config.Client(
			oauth2.NoContext,
			&oauth2.Token{
				AccessToken: user.Hash,
			},
		)
		client.SetClient(auther)

		// 2.c activate the repository
		_, err := client.RepoEnable(repo.Namespace, repo.Name)
		if err != nil {
			log.WithError(err).Errorf("activation failed")
			multierror.Append(result, err)
			continue
		}

		log.Debugln("successfully activated")
	}

	logrus.Infoln("repository activation complete")
	return result
}
