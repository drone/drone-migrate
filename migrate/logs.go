package migrate

import (
	"database/sql"
	"fmt"

	"github.com/russross/meddler"
	"github.com/sirupsen/logrus"
)

// MigrateLogs migrates the steps from the V0
// database to the V1 database.
func MigrateLogs(source, target *sql.DB) error {
	stepsV0 := []*StepV0{}

	// 1. load all stages from the V0 database.
	err := meddler.QueryAll(source, &stepsV0, "select * from procs where proc_ppid != 0")
	if err != nil {
		return err
	}

	logrus.Infof("migrating %d logs", len(stepsV0))

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
		logsV0 := &LogsV0{}
		err := meddler.QueryRow(source, logsV0, fmt.Sprintf("select * from logs where log_job_id = %d", stepV0.ID))
		if err != nil {
			logrus.WithError(err).Warnf("cannot find logs for step: id: %d", stepV0.ID)
			continue
		}

		logsV1 := &LogsV1{
			ID:   logsV0.ProcID,
			Data: logsV0.Data,
		}

		err = meddler.Insert(tx, "logs", logsV1)
		if err != nil {
			logrus.WithError(err).Errorln("migration failed")
			return err
		}
	}

	logrus.Infof("migration complete")
	return tx.Commit()
}
