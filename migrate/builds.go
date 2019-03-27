package migrate

import (
	"database/sql"
	"fmt"

	"github.com/russross/meddler"
	"github.com/sirupsen/logrus"
)

// MigrateBuilds migrates the builds from the V0
// database to the V1 database.
func MigrateBuilds(source, target *sql.DB) error {
	buildsV0 := []*BuildV0{}

	// 1. load all repos from the V0 database.
	err := meddler.QueryAll(source, &buildsV0, "select * from builds")
	if err != nil {
		return err
	}

	logrus.Infof("migrating %d builds", len(buildsV0))

	// 2. create a database transaction so that we
	// can rollback if the data migration fails.
	tx, err := target.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 3. iterate through the list and convert from
	// the 0.x to the 1.x structure and insert.
	var sequence int64
	for _, buildV0 := range buildsV0 {
		if buildV0.ID > sequence {
			sequence = buildV0.ID
		}

		log := logrus.
			WithField("repository", buildV0.RepoID).
			WithField("build", buildV0.Number)
		log.Debugln("migrate build")

		buildV1 := &BuildV1{
			ID:           buildV0.ID,
			RepoID:       buildV0.RepoID,
			Trigger:      "@hook",
			Number:       buildV0.Number,
			Parent:       buildV0.Parent,
			Status:       buildV0.Status,
			Error:        buildV0.Error,
			Event:        buildV0.Event,
			Action:       "",
			Link:         buildV0.Link,
			Timestamp:    buildV0.Timestamp,
			Title:        buildV0.Title,
			Message:      buildV0.Message,
			Before:       buildV0.Commit,
			After:        buildV0.Commit,
			Ref:          buildV0.Ref,
			Fork:         "",
			Source:       buildV0.Branch,
			Target:       buildV0.Branch,
			Author:       buildV0.Author,
			AuthorName:   buildV0.Author,
			AuthorEmail:  buildV0.Email,
			AuthorAvatar: buildV0.Avatar,
			Sender:       buildV0.Sender,
			Params:       map[string]string{},
			Deploy:       buildV0.Deploy,
			Started:      buildV0.Started,
			Finished:     buildV0.Finished,
			Created:      buildV0.Created,
			Updated:      buildV0.Created,
			Version:      1,
		}
		if len(buildV1.Message) > 1000 {
			buildV1.Message = buildV1.Message[:1000]
		}
		if len(buildV1.Title) > 1000 {
			buildV1.Title = buildV1.Title[:1000]
		}

		err = meddler.Insert(tx, "builds", buildV1)
		if err != nil {
			log.WithError(err).Errorln("migration failed")
			return err
		}

		//
		// migrate stages.
		//

		log.Debugln("build migration complete")
	}

	if meddler.Default == meddler.PostgreSQL {
		_, err = tx.Exec(fmt.Sprintf(updateBuildSeq, sequence))
		if err != nil {
			logrus.WithError(err).Errorln("failed to reset sequence")
			return err
		}
	}

	logrus.Infof("migration complete")
	return tx.Commit()
}

const buildListQuery = `
SELECT builds.*
FROM builds INNER JOIN repos ON build.build_repo_id = repos.repo_id
WHERE repos.repo_user_id > 0
`

const updateBuildSeq = `
ALTER SEQUENCE builds_build_id_seq
RESTART WITH %d
`
