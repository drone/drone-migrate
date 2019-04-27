package migrate

import (
	"database/sql"
	"fmt"

	"github.com/russross/meddler"
	"github.com/sirupsen/logrus"
)

// MigrateStages migrates the stages from the V0
// database to the V1 database.
func MigrateStages(source, target *sql.DB) error {
	stagesV0 := []*StageV0{}

	// 1. load all repos from the V0 database.
	err := meddler.QueryAll(source, &stagesV0, stageListQuery)
	if err != nil {
		return err
	}

	logrus.Infof("migrating %d stages", len(stagesV0))

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
	for _, stageV0 := range stagesV0 {
		if stageV0.ID > sequence {
			sequence = stageV0.ID
		}

		stageV1 := &StageV1{
			ID:        stageV0.ID,
			RepoID:    0,
			BuildID:   stageV0.BuildID,
			Number:    stageV0.PID,
			Name:      stageV0.Name,
			Kind:      "",
			Type:      "",
			Status:    stageV0.State,
			Error:     stageV0.Error,
			ErrIgnore: false,
			ExitCode:  stageV0.ExitCode,
			Machine:   stageV0.Machine,
			OS:        "linux",
			Arch:      "amd64",
			Variant:   "",
			Kernel:    "",
			Limit:     0,
			Started:   stageV0.Started,
			Stopped:   stageV0.Stopped,
			Created:   stageV0.Started,
			Updated:   stageV0.Stopped,
			Version:   1,
			OnSuccess: true,
			OnFailure: false,
			DependsOn: []string{},
			Labels:    map[string]string{},
		}
		if stageV1.Name == "" {
			stageV1.Name = "default"
		}

		err = meddler.Insert(tx, "stages", stageV1)
		if err != nil {
			logrus.WithError(err).Errorln("migration failed")
			return err
		}
	}

	if meddler.Default == meddler.PostgreSQL {
		_, err = tx.Exec(fmt.Sprintf(updateStageSeq, sequence+1))
		if err != nil {
			logrus.WithError(err).Errorln("failed to reset sequence")
			return err
		}
	}

	logrus.Infof("migration complete")
	return tx.Commit()
}

const stageListQuery = `
SELECT procs.*
FROM procs
INNER JOIN builds ON procs.proc_build_id = builds.build_id
INNER JOIN repos ON builds.build_repo_id = repos.repo_id
WHERE proc_ppid = 0
  AND repo_user_id > 0
`

const updateStageSeq = `
ALTER SEQUENCE stages_stage_id_seq
RESTART WITH %d
`
