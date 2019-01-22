package main

import (
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-migrate/migrate"
	"github.com/drone/drone-migrate/migrate/db"
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitea"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/drone/go-scm/scm/driver/gitlab"
	"github.com/drone/go-scm/scm/driver/gogs"
	"github.com/drone/go-scm/scm/driver/stash"
	"github.com/drone/go-scm/scm/transport/oauth1"
	"github.com/drone/go-scm/scm/transport/oauth2"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

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
		cli.StringFlag{
			Name:   "scm-driver",
			Usage:  "source code management system driver (github,gitlab,gogs,gitea,bitbucket,stash)",
			EnvVar: "SCM_DRIVER",
		},
		cli.StringFlag{
			Name:   "scm-server",
			Usage:  "source code management server address",
			EnvVar: "SCM_SERVER",
		},
		cli.StringFlag{
			Name:   "stash-consumer-key",
			Usage:  "atlassian stash consumer key",
			EnvVar: "STASH_CONSUMER_KEY",
		},
		cli.StringFlag{
			Name:   "stash-private-key-file",
			Usage:  "atlassian stash private key file",
			EnvVar: "STASH_PRIVATE_KEY_FILE",
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
				var (
					driver     = c.GlobalString("target-database-driver")
					datasource = c.GlobalString("target-database-datasource")
					provider   = c.GlobalString("scm-driver")
					server     = c.GlobalString("scm-server")
				)

				logrus.Debugf("target database driver: %s", driver)
				logrus.Debugf("target database datasource: %s", datasource)
				logrus.Debugf("scm driver: %s", provider)
				logrus.Debugf("scm server: %s", server)

				target, err := sql.Open(driver, datasource)

				if err != nil {
					return err
				}

				client, err := createClient(c)

				if err != nil {
					return err
				}

				return migrate.UpdateRepoIdentifiers(target, client)
			},
		},
		{
			Name:  "migrate-secrets",
			Usage: "migrate drone secrets",
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

				return migrate.MigrateSecrets(source, target)
			},
		},
		{
			Name:  "migrate-registries",
			Usage: "migrate registry credentials",
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

				return migrate.MigrateRegistries(source, target)
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

				return migrate.ActivateRepositories(
					target,
					drone.New(c.GlobalString("drone-server")),
				)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func createClient(c *cli.Context) (*scm.Client, error) {
	server := c.GlobalString("scm-server")

	switch c.GlobalString("scm-driver") {
	case "gogs":
		client, err := gogs.New(server)

		if err != nil {
			return nil, err
		}

		client.Client = &http.Client{
			Transport: &oauth2.Transport{
				Scheme: oauth2.SchemeToken,
				Source: oauth2.ContextTokenSource(),
			},
		}

		return client, nil
	case "gitea":
		client, err := gitea.New(server)

		if err != nil {
			return nil, err
		}

		client.Client = &http.Client{
			Transport: &oauth2.Transport{
				Scheme: oauth2.SchemeToken,
				Source: oauth2.ContextTokenSource(),
			},
		}

		return client, nil
	case "gitlab":
		client, err := gitlab.New(server)

		if err != nil {
			return nil, err
		}

		client.Client = &http.Client{
			Transport: &oauth2.Transport{
				Source: oauth2.ContextTokenSource(),
			},
		}

		return client, nil
	case "github":
		client, err := github.New(server)

		if err != nil {
			return nil, err
		}

		client.Client = &http.Client{
			Transport: &oauth2.Transport{
				Source: oauth2.ContextTokenSource(),
			},
		}

		return client, nil
	case "stash":
		privateKey, err := parsePrivateKeyFile(
			c.String("stash-private-key-file"),
		)

		if err != nil {
			return nil, err
		}

		client, err := stash.New(server)

		if err != nil {
			return nil, err
		}

		client.Client = &http.Client{
			Transport: &oauth1.Transport{
				ConsumerKey: c.String("stash-consumer-key"),
				PrivateKey:  privateKey,
				Source:      oauth1.ContextTokenSource(),
			},
		}
		return client, nil
	default:
		return nil, errors.New("Source code management system not configured")
	}
}

// parsePrivateKeyFile is a helper function that parses an
// RSA Private Key file encoded in PEM format.
func parsePrivateKeyFile(path string) (*rsa.PrivateKey, error) {
	d, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return parsePrivateKey(d)
}

// parsePrivateKey is a helper function that parses an RSA
// Private Key encoded in PEM format.
func parsePrivateKey(data []byte) (*rsa.PrivateKey, error) {
	p, _ := pem.Decode(data)
	return x509.ParsePKCS1PrivateKey(p.Bytes)
}
