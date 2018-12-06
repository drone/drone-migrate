package main

import (
	"database/sql"
	"os"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-migrate/migrate"
	"github.com/drone/drone-migrate/migrate/db"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// TODO(bradrydzewski) update repository secrets
// TODO(bradrydzewski) update builds
// TODO(bradrydzewski) update stages
// TODO(bradrydzewski) update steps
// TODO(bradrydzewski) update logs

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "source-database-driver",
			Usage:  "Source database driver",
			EnvVar: "SOURCE_DATABASE_DRIVER",
		},
		cli.StringFlag{
			Name:   "source-database-datasource",
			Usage:  "Source database datasource",
			EnvVar: "SOURCE_DATABASE_DATASOURCE",
		},
		cli.StringFlag{
			Name:   "target-database-driver",
			Usage:  "target database driver",
			EnvVar: "TARGET_DATABASE_DRIVER",
		},
		cli.StringFlag{
			Name:   "target-database-datasource",
			Usage:  "target database datasource",
			EnvVar: "TARGET_DATABASE_DATASOURCE",
		},
		cli.StringFlag{
			Name:   "drone-server",
			Usage:  "target drone server address",
			EnvVar: "DRONE_SERVER",
		},
		cli.BoolTFlag{
			Name:   "debug",
			Usage:  "enable debug mode",
			EnvVar: "DEBUG",
		},
	}

	app.Before = func(c *cli.Context) error {
		if c.GlobalBoolT("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "setup-database",
			Usage: "initialize the 1.0 database",
			Action: func(c *cli.Context) error {
				var (
					driver     = c.GlobalString("target-database-driver")
					datasource = c.GlobalString("target-database-datasource")
				)
				logrus.Debugf("target database driver: %s", driver)
				logrus.Debugf("target database datasource: %s", datasource)
				target, err := sql.Open(driver, datasource)
				if err != nil {
					return err
				}
				err = db.Create(target, driver)
				if err != nil {
					return err
				}
				logrus.Infoln("target database created")
				return nil
			},
		},
		{
			Name:  "migrate-users",
			Usage: "migrate user resources",
			Action: func(c *cli.Context) error {
				source, err := sql.Open(
					c.GlobalString("source-database-driver"),
					c.GlobalString("source-database-datasource"),
				)
				if err != nil {
					return err
				}
				target, err := sql.Open(
					c.GlobalString("target-database-driver"),
					c.GlobalString("target-database-datasource"),
				)
				if err != nil {
					return err
				}
				return migrate.MigrateUsers(source, target)
			},
		},
		{
			Name:  "migrate-repos",
			Usage: "migrate repository resources",
			Action: func(c *cli.Context) error {
				source, err := sql.Open(
					c.GlobalString("source-database-driver"),
					c.GlobalString("source-database-datasource"),
				)
				if err != nil {
					return err
				}
				target, err := sql.Open(
					c.GlobalString("target-database-driver"),
					c.GlobalString("target-database-datasource"),
				)
				if err != nil {
					return err
				}
				return migrate.MigrateRepos(source, target)
			},
		},
		{
			Name:  "update-repos",
			Usage: "update repository metadata",
			Action: func(c *cli.Context) error {
				target, err := sql.Open(
					c.GlobalString("target-database-driver"),
					c.GlobalString("target-database-datasource"),
				)
				if err != nil {
					return err
				}
				//
				// TODO create the remote repository client.
				//
				return migrate.UpdateRepoIdentifiers(target, nil)
			},
		},
		{
			Name:  "activate-repos",
			Usage: "activate repository resources",
			Action: func(c *cli.Context) error {
				target, err := sql.Open(
					c.GlobalString("target-database-driver"),
					c.GlobalString("target-database-datasource"),
				)
				if err != nil {
					return err
				}
				client := drone.New(c.GlobalString("drone-server"))
				return migrate.ActivateRepositories(target, client)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
