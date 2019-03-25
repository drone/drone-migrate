package migrate

import (
	"database/sql"
	"fmt"

	"github.com/russross/meddler"
	"github.com/sirupsen/logrus"
)

// MigrateSteps migrates the steps from the V0
// database to the V1 database.
func MigrateSteps(source, target *sql.DB) error {
	stepsV0 := []*StepV0{}

	// 1. load all stages from the V0 database.
	err := meddler.QueryAll(source, &stepsV0, stepListQuery)
	if err != nil {
		return err
	}

	logrus.Infof("migrating %d steps", len(stepsV0))

	// 2. create a database transaction so that we
	// can rollback if the data migration fails.
	tx, err := target.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 3. iterate through the list and convert from
	// the 0.x to the 1.x structure and insert.
	for _, stepV0 := range stepsV0 {
		stageV0 := &StageV0{}
		err := meddler.QueryRow(source, stageV0, fmt.Sprintf("select * from procs where proc_pid = %d and proc_build_id = %d", stepV0.PPID, stepV0.BuildID))
		if err != nil {
			logrus.WithError(err).Errorln("cannot find parent step")
			return err
		}

		stepV1 := &StepV1{
			ID:        stepV0.ID,
			StageID:   stageV0.ID,
			Number:    stepV0.PID,
			Name:      stepV0.Name,
			Status:    stepV0.State,
			Error:     stepV0.Error,
			ErrIgnore: false,
			ExitCode:  stepV0.ExitCode,
			Started:   stepV0.Started,
			Stopped:   stepV0.Stopped,
			Version:   1,
		}

		err = meddler.Insert(tx, "steps", stepV1)
		if err != nil {
			logrus.WithError(err).Errorln("migration failed")
			return err
		}
	}

	logrus.Infof("migration complete")
	return tx.Commit()
}

const stepListQuery = `
SELECT procs.*
FROM procs
INNER JOIN builds ON procs.proc_build_id = builds.build_id
INNER JOIN repos ON builds.build_repo_id = repos.repo_id
WHERE proc_ppid != 0
  AND repo_user_id > 0
`
