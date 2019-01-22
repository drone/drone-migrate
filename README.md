__IMPORTANT: this migration utility is incomplete. It will not yet generate a working 1.0 database__

Migration utility will help you migrate from a Drone 0.8.x database to a Drone 1.0.x database.

## Preparing for the migration

1. create a full backup of your 0.8.x Drone database
2. create a new database for your 1.0.x server
3. do not login to your 1.0 drone instance until this is complete

## Building the migration utility

```shell
go get -u github.com/drone/drone-migrate
```

## Configuring the migration utility

The migration utility will copy data from your 0.8.x database to your new 1.0.x database. You will need to provide the migration tool with the connection string for both the old and new database.

```sh
export SOURCE_DATABASE_DRIVER=sqlite3|mysql|postgres
export TARGET_DATABASE_DRIVER=sqlite3|mysql|postgres
export SOURCE_DATABASE_DATASOURCE=/path/to/old/database.sqlite
export TARGET_DATABASE_DATASOURCE=/path/to/new/database.sqlite
```

If you are using GitHub, configure the GitHub driver:

```sh
export SCM_DRIVER=github
export SCM_SERVER=https://api.github.com
```

If you are using GitHub Enterprise, configure the GitHub driver:

```sh
export SCM_DRIVER=github
export SCM_SERVER=https://github.company.com/api/v3
```

If you are using Gogs, configure the Gogs driver:

```sh
export SCM_DRIVER=gogs
export SCM_SERVER=https://gogs.company.com
```

If you are using Gitea, configure the Gitea driver:

```sh
export SCM_DRIVER=gitea
export SCM_SERVER=https://gitea.company.com
```

If you are using Stash, configure the Stash driver:

```sh
export SCM_DRIVER=stash
export SCM_SERVER=https://stash.company.com
export STASH_CONSUMER_KEY=OauthKey
export STASH_PRIVATE_KEY_FILE=/path/to/private/key.pem
```

## Create the 1.0 database

```shell
$ drone-migrate setup-database
```

## Migrate users from 0.8 to 1.0

```shell
$ drone-migrate migrate-users
```

## Migrate repositories from 0.8 to 1.0

```shell
$ drone-migrate migrate-repos
```

## Update the repository metadata

Drone 1.0 stores addition repository metadata that needs to be fetched from the source code management system. This additional metadata is required.

```shell
$ drone-migrate update-repos
```

## Migrate secrets from 0.8 to 1.0

Secrets stored within Drone can be migrated, if you use some external tool to store your secrets like Vault you can skip this step.

```shell
$ drone-migrate migrate-secrets
```

## Migrate registry credentials from 0.8 to 1.0

If you haven't used ayn private images within the pipeline you can skip this step, this is only needed if you are using private Docker images for your Drone steps.

```shell
$ drone-migrate migrate-registires
```

## Activate the repositories

The final step is to ensure all repositories are activated and have a valid web-hook configured in the source code management system.

```shell
$ drone-migrate activate-repos
```
