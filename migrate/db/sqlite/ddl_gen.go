package sqlite

import (
	"database/sql"
)

var migrations = []struct {
	name string
	stmt string
}{
	{
		name: "create-table-users",
		stmt: createTableUsers,
	},
	{
		name: "create-table-repos",
		stmt: createTableRepos,
	},
	{
		name: "create-table-perms",
		stmt: createTablePerms,
	},
	{
		name: "create-index-perms-user",
		stmt: createIndexPermsUser,
	},
	{
		name: "create-index-perms-repo",
		stmt: createIndexPermsRepo,
	},
	{
		name: "create-table-builds",
		stmt: createTableBuilds,
	},
	{
		name: "create-index-builds-in-progress",
		stmt: createIndexBuildsInProgress,
	},
	{
		name: "create-index-builds-repo",
		stmt: createIndexBuildsRepo,
	},
	{
		name: "create-index-builds-author",
		stmt: createIndexBuildsAuthor,
	},
	{
		name: "create-index-builds-sender",
		stmt: createIndexBuildsSender,
	},
	{
		name: "create-index-builds-ref",
		stmt: createIndexBuildsRef,
	},
	{
		name: "create-index-build-incomplete",
		stmt: createIndexBuildIncomplete,
	},
	{
		name: "create-table-stages",
		stmt: createTableStages,
	},
	{
		name: "create-index-stages-build",
		stmt: createIndexStagesBuild,
	},
	{
		name: "create-index-stages-status",
		stmt: createIndexStagesStatus,
	},
	{
		name: "create-table-steps",
		stmt: createTableSteps,
	},
	{
		name: "create-index-steps-stage",
		stmt: createIndexStepsStage,
	},
	{
		name: "create-table-logs",
		stmt: createTableLogs,
	},
	{
		name: "create-table-cron",
		stmt: createTableCron,
	},
	{
		name: "create-index-cron-repo",
		stmt: createIndexCronRepo,
	},
	{
		name: "create-index-cron-next",
		stmt: createIndexCronNext,
	},
	{
		name: "create-table-secrets",
		stmt: createTableSecrets,
	},
	{
		name: "create-index-secrets-repo",
		stmt: createIndexSecretsRepo,
	},
	{
		name: "create-index-secrets-repo-name",
		stmt: createIndexSecretsRepoName,
	},
	{
		name: "create-table-nodes",
		stmt: createTableNodes,
	},
}

// Migrate performs the database migration. If the migration fails
// and error is returned.
func Migrate(db *sql.DB) error {
	if err := createTable(db); err != nil {
		return err
	}
	completed, err := selectCompleted(db)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for _, migration := range migrations {
		if _, ok := completed[migration.name]; ok {

			continue
		}

		if _, err := db.Exec(migration.stmt); err != nil {
			return err
		}
		if err := insertMigration(db, migration.name); err != nil {
			return err
		}

	}
	return nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(migrationTableCreate)
	return err
}

func insertMigration(db *sql.DB, name string) error {
	_, err := db.Exec(migrationInsert, name)
	return err
}

func selectCompleted(db *sql.DB) (map[string]struct{}, error) {
	migrations := map[string]struct{}{}
	rows, err := db.Query(migrationSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		migrations[name] = struct{}{}
	}
	return migrations, nil
}

//
// migration table ddl and sql
//

var migrationTableCreate = `
CREATE TABLE IF NOT EXISTS migrations (
 name VARCHAR(255)
,UNIQUE(name)
)
`

var migrationInsert = `
INSERT INTO migrations (name) VALUES (?)
`

var migrationSelect = `
SELECT name FROM migrations
`

//
// 001_create_table_user.sql
//

var createTableUsers = `
CREATE TABLE IF NOT EXISTS users (
 user_id            INTEGER PRIMARY KEY AUTOINCREMENT
,user_login         TEXT
,user_email         TEXT
,user_admin         BOOLEAN
,user_machine       BOOLEAN
,user_active        BOOLEAN
,user_avatar        TEXT
,user_syncing       BOOLEAN
,user_synced        INTEGER
,user_created       INTEGER
,user_updated       INTEGER
,user_last_login    INTEGER
,user_oauth_token   TEXT
,user_oauth_refresh TEXT
,user_oauth_expiry  INTEGER
,user_hash          TEXT
,UNIQUE(user_login COLLATE NOCASE)
,UNIQUE(user_hash)
);
`

//
// 002_create_table_repos.sql
//

var createTableRepos = `
CREATE TABLE IF NOT EXISTS repos (
 repo_id                    INTEGER PRIMARY KEY AUTOINCREMENT
,repo_uid                   TEXT
,repo_user_id               INTEGER
,repo_namespace             TEXT
,repo_name                  TEXT
,repo_slug                  TEXT
,repo_scm                   TEXT
,repo_clone_url             TEXT
,repo_ssh_url               TEXT
,repo_html_url              TEXT
,repo_active                BOOLEAN
,repo_private               BOOLEAN
,repo_visibility            TEXT
,repo_branch                TEXT
,repo_counter               INTEGER
,repo_config                TEXT
,repo_timeout               INTEGER
,repo_trusted               BOOLEAN
,repo_protected             BOOLEAN
,repo_synced                INTEGER
,repo_created               INTEGER
,repo_updated               INTEGER
,repo_version               INTEGER
,repo_signer                TEXT
,repo_secret                TEXT
,UNIQUE(repo_slug)
,UNIQUE(repo_uid)
);
`

//
// 003_create_table_perms.sql
//

var createTablePerms = `
CREATE TABLE IF NOT EXISTS perms (
 perm_user_id  INTEGER
,perm_repo_uid TEXT
,perm_read     BOOLEAN
,perm_write    BOOLEAN
,perm_admin    BOOLEAN
,perm_synced   INTEGER
,perm_created  INTEGER
,perm_updated  INTEGER
,PRIMARY KEY(perm_user_id, perm_repo_uid)
);
`

var createIndexPermsUser = `
CREATE INDEX IF NOT EXISTS ix_perms_user ON perms (perm_user_id);
`

var createIndexPermsRepo = `
CREATE INDEX IF NOT EXISTS ix_perms_repo ON perms (perm_repo_uid);
`

//
// 004_create_table_builds.sql
//

var createTableBuilds = `
CREATE TABLE IF NOT EXISTS builds (
 build_id            INTEGER PRIMARY KEY AUTOINCREMENT
,build_repo_id       INTEGER
,build_trigger       TEXT
,build_number        INTEGER
,build_parent        INTEGER
,build_status        TEXT
,build_error         TEXT
,build_event         TEXT
,build_action        TEXT
,build_link          TEXT
,build_timestamp     INTEGER
,build_title         TEXT
,build_message       TEXT
,build_before        TEXT
,build_after         TEXT
,build_ref           TEXT
,build_source_repo   TEXT
,build_source        TEXT
,build_target        TEXT
,build_author        TEXT
,build_author_name   TEXT
,build_author_email  TEXT
,build_author_avatar TEXT
,build_sender        TEXT
,build_deploy        TEXT
,build_params        TEXT
,build_started       INTEGER
,build_finished      INTEGER
,build_created       INTEGER
,build_updated       INTEGER
,build_version       INTEGER
,UNIQUE(build_repo_id, build_number)
);
`

var createIndexBuildsInProgress = `
CREATE INDEX IF NOT EXISTS ix_build_in_progress ON builds (build_status)
WHERE build_status IN ('pending', 'running');
`

var createIndexBuildsRepo = `
CREATE INDEX IF NOT EXISTS ix_build_repo ON builds (build_repo_id);
`

var createIndexBuildsAuthor = `
CREATE INDEX IF NOT EXISTS ix_build_author ON builds (build_author);
`

var createIndexBuildsSender = `
CREATE INDEX IF NOT EXISTS ix_build_sender ON builds (build_sender);
`

var createIndexBuildsRef = `
CREATE INDEX IF NOT EXISTS ix_build_ref ON builds (build_repo_id, build_ref);
`

var createIndexBuildIncomplete = `
CREATE INDEX IF NOT EXISTS ix_build_incomplete ON builds (build_status)
WHERE build_status IN ('pending', 'running');
`

//
// 005_create_table_stages.sql
//

var createTableStages = `
CREATE TABLE IF NOT EXISTS stages (
 stage_id          INTEGER PRIMARY KEY AUTOINCREMENT
,stage_repo_id     INTEGER
,stage_build_id    INTEGER
,stage_number      INTEGER
,stage_kind        TEXT
,stage_type        TEXT
,stage_name        TEXT
,stage_status      TEXT
,stage_error       TEXT
,stage_errignore   BOOLEAN
,stage_exit_code   INTEGER
,stage_limit       INTEGER
,stage_os          TEXT
,stage_arch        TEXT
,stage_variant     TEXT
,stage_kernel      TEXT
,stage_machine     TEXT
,stage_started     INTEGER
,stage_stopped     INTEGER
,stage_created     INTEGER
,stage_updated     INTEGER
,stage_version     INTEGER
,stage_on_success  BOOLEAN
,stage_on_failure  BOOLEAN
,stage_depends_on  TEXT
,stage_labels      TEXT
,UNIQUE(stage_build_id, stage_number)
,FOREIGN KEY(stage_build_id) REFERENCES builds(build_id) ON DELETE CASCADE
);
`

var createIndexStagesBuild = `
CREATE INDEX IF NOT EXISTS ix_stages_build ON stages (stage_build_id);
`

var createIndexStagesStatus = `
CREATE INDEX IF NOT EXISTS ix_build_in_progress ON stages (stage_status)
WHERE stage_status IN ('pending', 'running');
`

//
// 006_create_table_steps.sql
//

var createTableSteps = `
CREATE TABLE IF NOT EXISTS steps (
 step_id          INTEGER PRIMARY KEY AUTOINCREMENT
,step_stage_id    INTEGER
,step_number      INTEGER
,step_name        TEXT
,step_status      TEXT
,step_error       TEXT
,step_errignore   BOOLEAN
,step_exit_code   INTEGER
,step_started     INTEGER
,step_stopped     INTEGER
,step_version     INTEGER
,UNIQUE(step_stage_id, step_number)
,FOREIGN KEY(step_stage_id) REFERENCES stages(stage_id) ON DELETE CASCADE
);
`

var createIndexStepsStage = `
CREATE INDEX IF NOT EXISTS ix_steps_stage ON steps (step_stage_id);
`

//
// 007_create_table_logs.sql
//

var createTableLogs = `
CREATE TABLE IF NOT EXISTS logs (
 log_id    INTEGER PRIMARY KEY
,log_data  BLOB
,FOREIGN KEY(log_id) REFERENCES steps(step_id) ON DELETE CASCADE
);
`

//
// 008_create_table_cron.sql
//

var createTableCron = `
CREATE TABLE IF NOT EXISTS cron (
 cron_id          INTEGER PRIMARY KEY AUTOINCREMENT
,cron_repo_id     INTEGER
,cron_name        TEXT
,cron_expr        TEXT
,cron_next        INTEGER
,cron_prev        INTEGER
,cron_event       TEXT
,cron_branch      TEXT
,cron_target      TEXT
,cron_disabled    BOOLEAN
,cron_created     INTEGER
,cron_updated     INTEGER
,cron_version     INTEGER
,UNIQUE(cron_repo_id, cron_name)
,FOREIGN KEY(cron_repo_id) REFERENCES repos(repo_id) ON DELETE CASCADE
);
`

var createIndexCronRepo = `
CREATE INDEX IF NOT EXISTS ix_cron_repo ON cron (cron_repo_id);
`

var createIndexCronNext = `
CREATE INDEX IF NOT EXISTS ix_cron_next ON cron (cron_next);
`

//
// 009_create_table_secrets.sql
//

var createTableSecrets = `
CREATE TABLE IF NOT EXISTS secrets (
 secret_id                INTEGER PRIMARY KEY AUTOINCREMENT
,secret_repo_id           INTEGER
,secret_name              TEXT
,secret_data              BLOB
,secret_pull_request      BOOLEAN
,secret_pull_request_push BOOLEAN
,UNIQUE(secret_repo_id, secret_name)
,FOREIGN KEY(secret_repo_id) REFERENCES repos(repo_id) ON DELETE CASCADE
);
`

var createIndexSecretsRepo = `
CREATE INDEX IF NOT EXISTS ix_secret_repo ON secrets (secret_repo_id);
`

var createIndexSecretsRepoName = `
CREATE INDEX IF NOT EXISTS ix_secret_repo_name ON secrets (secret_repo_id, secret_name);
`

//
// 010_create_table_nodes.sql
//

var createTableNodes = `
CREATE TABLE IF NOT EXISTS nodes (
 node_id         INTEGER PRIMARY KEY AUTOINCREMENT
,node_uid        TEXT
,node_provider   TEXT
,node_state      TEXT
,node_name       TEXT
,node_image      TEXT
,node_region     TEXT
,node_size       TEXT
,node_os         TEXT
,node_arch       TEXT
,node_kernel     TEXT
,node_variant    TEXT
,node_address    TEXT
,node_capacity   INTEGER
,node_filter     TEXT
,node_labels     TEXT
,node_error      TEXT
,node_ca_key     TEXT
,node_ca_cert    TEXT
,node_tls_key    TEXT
,node_tls_cert   TEXT
,node_tls_name   TEXT
,node_paused     BOOLEAN
,node_protected  BOOLEAN
,node_created    INTEGER
,node_updated    INTEGER
,node_pulled     INTEGER

,UNIQUE(node_name)
);
`
