Migration utility will help you migrate from a Drone 0.8.x database to a Drone 1.0.x database.

___Please note the migration utility may require manual database cleanup.___ For example, in 0.8 the same repository can be listed in the database twice if it has been renamed, however, in 1.0 this will cause unique key violations. These edge cases require manual intervention. You should therefore be comfortable with sql and database troubleshooting before you proceed.

## Preparing for the migration

1. create a full backup of your 0.8.x Drone database
2. create a new database for your 1.0.x server
3. do not create or start your drone 1.0 container until this is complete

## Download the migration utility

```
docker pull drone/migrate
```

## Configuring the migration utility

The migration utility will copy data from your 0.8.x database to your new _empty_ 1.0.x database. You will need to provide the migration tool with the connection string for both the old and new database.

```sh
-e SOURCE_DATABASE_DRIVER=sqlite3|mysql|postgres
-e TARGET_DATABASE_DRIVER=sqlite3|mysql|postgres
-e SOURCE_DATABASE_DATASOURCE=/path/to/old/database.sqlite|user:password@tcp(hostname:port)/database
-e TARGET_DATABASE_DATASOURCE=/path/to/new/database.sqlite|user:password@tcp(hostname:port)/database
```

Configure the Drone 1.0 server address:

```
-e DRONE_SERVER=https://drone.company.com
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

If you are using Bitbucket Cloud, configure the Bitbucket driver:

```sh
-e SCM_DRIVER=bitbucket
-e BITBUCKET_CLIENT_ID=$your_client_id
-e BITBUCKET_CLIENT_SECRET=$your_client_secret
```

# Full Migration

Run the below migration steps to copy your data from your 0.8 database to your 1.x database. _Running the below commands will not have any impact on your existing 0.8 database._

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
$ docker run -e [...] drone/migrate remove-renamed
$ docker run -e [...] drone/migrate remove-not-found
```

## Optional Encryption

You can also optionally [configure](https://docs.drone.io/server/storage/encryption/) secret encryption in Drone 1.0. If you plan on enabling encryption you will need to encrypt the secrets before you complete the migration.

```
$ export TARGET_DATABASE_ENCRYPTION_KEY=....
$ docker run -e [...] -e drone-drone/migrate encrypt-secrets
```

## Final Migration Step

The final step is to re-activate your repositories. At this time it is safe to start your Drone server. Once the server is started you can execute the final migration command:

```
$ docker run -e [...] drone/migrate activate-repos
```

_The above command should only be executed once you are ready to finalize your migration.  When you execute this command it may replace any webhooks created by your 0.8 instance with webhooks that point to your 1.x instance._


## Final Check

The migration is not considered complete until all steps are completed and the below sql query returns an empty result set.  If the below query does not return an empty result set you should execute the optional `remove-renamed` and `remove-not-found` migration steps.

```text
SELECT *
FROM repos
WHERE repo_uid LIKE 'temp_%'
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

you can optionally migrate logs to s3 storage. _Note that the migration utility authenticates with aws using standard authentication methods, including aws_access_key_id and aws_secret_access_key_


```shell
$ docker run -e S3_BUCKET=<bucket> -e [...] drone/migrate migrate-logs-s3
```

## Migrate secrets from 0.8 to 1.0

Secrets stored within Drone can be migrated, if you use some external tool to store your secrets like Vault you can skip this step.

```shell
$ docker run -e [...] drone/migrate migrate-secrets
```

## Migrate registry credentials from 0.8 to 1.0

If you haven't used ayn private images within the pipeline you can skip this step, this is only needed if you are using private Docker images for your Drone steps.

```shell
$ docker run -e [...] drone/migrate migrate-registries
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

## Dump Tokens (Optional)

You can optionally dump 0.8 user API tokens for use with 1.0 as described [here](https://github.com/drone/drone/issues/2713). If your team heavily uses Drone tokens in their build process (to trigger downstream builds, etc) you may find this helpful.

```shell
$ docker run -e [...] drone/migrate dump-tokens
```
