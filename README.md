__IMPORTANT: this migration utility is incomplete. It will not yet generate a working 1.0 database__

Migration utility will help you migrate from a Drone 0.8.x database to a Drone 1.0.x database.

## Preparing for the migration

1. create a full backup of your 0.8.x Drone database
2. create a new database for your 1.0.x server
3. do not create or start your drone 1.0 container until this is complete

## Download the migration utility

```
docker pull drone/migrate
```

## Configuring the migration utility

The migration utility will copy data from your 0.8.x database to your new 1.0.x database. You will need to provide the migration tool with the connection string for both the old and new database.

```sh
-e SOURCE_DATABASE_DRIVER=sqlite3|mysql|postgres
-e TARGET_DATABASE_DRIVER=sqlite3|mysql|postgres
-e SOURCE_DATABASE_DATASOURCE=/path/to/old/database.sqlite
-e TARGET_DATABASE_DATASOURCE=/path/to/new/database.sqlite
```

If you are using GitHub, configure the GitHub driver:

```sh
-e SCM_DRIVER=github
-e SCM_SERVER=https://api.github.com
```

If you are using GitHub Enterprise, configure the GitHub driver:

```sh
-e SCM_DRIVER=github
-e SCM_SERVER=https://github.company.com/api/v3
```

If you are using Gogs, configure the Gogs driver:

```sh
-e SCM_DRIVER=gogs
-e SCM_SERVER=https://gogs.company.com
```

If you are using Gitea, configure the Gitea driver:

```sh
-e SCM_DRIVER=gitea
-e SCM_SERVER=https://gitea.company.com
```

If you are using Stash, configure the Stash driver:

```sh
-e SCM_DRIVER=stash
-e SCM_SERVER=https://stash.company.com
-e STASH_CONSUMER_KEY=OauthKey
-e STASH_PRIVATE_KEY_FILE=/path/to/private/key.pem
```

# Full Migration


```
$ docker run -e [...] drone/migrate setup-database
$ docker run -e [...] drone/migrate migrate-users
$ docker run -e [...] drone/migrate migrate-repos
$ docker run -e [...] drone/migrate migrate-secrets
$ docker run -e [...] drone/migrate migrate-registries
$ docker run -e [...] drone/migrate migrate-builds
$ docker run -e [...] drone/migrate migrate-stages
$ docker run -e [...] drone/migrate migrate-steps
$ docker run -e [...] drone/migrate migrate-logs
$ docker run -e [...] drone/migrate update-repos
$ docker run -e [...] drone/migrate activate-repos
```

# Execution Individual Commands

This can be helpful if a particular migration step fails. You can safely truncate the impacted database table and then re-try the migration.

## Create the 1.0 database

```shell
$ docker run -e [...] drone/migrate setup-database
```

## Migrate users from 0.8 to 1.0

```shell
$ docker run -e [...] drone/migrate migrate-users
```

## Migrate repositories from 0.8 to 1.0

```shell
$ docker run -e [...] drone/migrate migrate-repos
```

## Migrate builds from 0.8 to 1.0

```shell
$ docker run -e [...] drone/migrate migrate-builds
```

## Migrate stages from 0.8 to 1.0

```shell
$ docker run -e [...] drone/migrate migrate-stages
```

## Migrate steps from 0.8 to 1.0

```shell
$ docker run -e [...] drone/migrate migrate-steps
```

## Migrate logs from 0.8 to 1.0

```shell
$ docker run -e [...] drone/migrate migrate-logs
```

## Migrate secrets from 0.8 to 1.0

Secrets stored within Drone can be migrated, if you use some external tool to store your secrets like Vault you can skip this step.

```shell
$ docker run -e [...] drone/migrate migrate-secrets
```

## Migrate registry credentials from 0.8 to 1.0

If you haven't used ayn private images within the pipeline you can skip this step, this is only needed if you are using private Docker images for your Drone steps.

```shell
$ docker run -e [...] drone/migrate migrate-registires
```

## Update the repository metadata

Drone 1.0 stores addition repository metadata that needs to be fetched from the source code management system. This additional metadata is required.

```shell
$ docker run -e [...] drone/migrate update-repos
```

## Activate the repositories

The final step is to ensure all repositories are activated and have a valid web-hook configured in the source code management system.

```shell
$ docker run -e [...] drone/migrate activate-repos
```
