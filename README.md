Migration utility will help you migrate from a Drone 0.8.x database to a Drone 1.0.x database. __IMPORTANT: this migration utility is incomplete. It will not yet generate a working 1.0 database__

# Preparing for the migration

1. create a full backup of your 0.8.x Drone database
2. create a new database for your 1.0.x server

# Building the migration utility

```shell
go get -u github.com/drone/drone-migrate
```

# Configuring the migration utility

The migration utility will copy data from your 0.8.x database to your new 1.0.x database. You will need to provide the migration tool with the connection string for both the old and new database.

```sh
export SOURCE_DATABASE_DRIVER=sqlite3|mysql|postgres
export TARGET_DATABASE_DRIVER=sqlite3|mysql|postgres
export SOURCE_DATABASE_DATASOURCE=/path/to/v0/database.sqlite
export TARGET_DATABASE_DATASOURCE=/path/to/v1/database.sqlite

# this should be set to your 1.0 server address.
# the address must omit the trailing slash.
export DRONE_SERVER=https://drone.company.com
```

# Create the 1.0 database

```shell
$ drone-migrate setup-database
```

# Migrate the user data from 0.8 to 1.0

```shell
$ drone-migrate migrate-users
```

# Migrate the repository data from 0.8 to 1.0

```shell
$ drone-migrate migrate-repos
```

# Migrate the repository secrets from 0.8 to 1.0

TODO

# Update the repository metadata

TODO

# Re-activate the repositories.

TODO



<!--
# Update the repository metadata

The latest version of Drone captures new fields that need to be retrieved from your source code management system (e.g. GitHub).

```shell
$ drone-migrate update-repos
```

# Activate the repositories.

The final step is to ensure all repositories are activated and have a valid web-hook configured in  the source code management system.

```shell
$ drone-migrate activate-repos
```
-->

<!--
NOTES:

1. we do not need to pass the drone token, we can get from the database
2. we do not need to pass the remote credentials, we can also get from the database
-->
