package migrate

import (
	"database/sql"

	"github.com/russross/meddler"
	"github.com/sirupsen/logrus"
)

// MigrateStages migrates the stages from the V0
// database to the V1 database.
func MigrateStages(source, target *sql.DB) error {
	stagesV0 := []*StageV0{}

	// 1. load all repos from the V0 database.
	err := meddler.QueryAll(source, &stagesV0, "select * from procs where proc_ppid = 0")
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
	for _, stageV0 := range stagesV0 {
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
			stageV1.Name = "defult"
		}

		err = meddler.Insert(tx, "stages", stageV1)
		if err != nil {
			logrus.WithError(err).Errorln("migration failed")
			return err
		}
	}

	logrus.Infof("migration complete")
	return tx.Commit()
}
