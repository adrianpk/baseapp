# Custom App

A base webapp based on [kabestan](https://gitlab.com/kabestan/backend/kabestan)

<img src="docs/img/user_index.png" width="480">

[More screenshots](docs/screenshots.md)

## Planned features

- Config
- Logging
- Authentication
- Authorization
- Localization (I10n)
- Web Endpoints
- JSON REST Endpoints
- gRPC API

## Changelog

## **20200128**

- Authorization middleware.

## **20200127**

- RBAC administration.

[Full changelog](docs/changelog.md)

## Routes

### Auth

| Method | Path          | Handler    |
| ------ | ------------- | ---------- |
| GET    | /auth/signup  | InitSignUp |
| POST   | /auth/signup  | SignUp     |
| GET    | /auth/signin  | InitSignIn |
| POST   | /auth/signin  | SignIn     |
| GET    | /auth/signout | SignOut    |

[More routes](docs/routes.md)

## Dev. Env. Setup [draft]

After stabilizing this base app a (still not published) generator will be updated to automate most of these steps.

## App Description

[More info](docs/description.md)

### Clone app

```shell
$ git clone https://gitlab.com/kabestan/repo/baseapp appname
```

Replace appname by the name of your app.

## Create database user

If it does not exist yet.

```shell
$ psql
psql (11.5 (Ubuntu 11.5-1))
Type "help" for help.

user=# CREATE ROLE rolename;
user=# ALTER USER rolename WITH PASSWORD 'password';
```

Replace rolename by the database user owner.
Replace password by prefered password.

**Note**: First migration step includes POSTGis extension installation.
To do so, the superuser permission is required.
You can use something like this instead to create a role with SUPERUSER power.

```shell
user=# CREATE ROLE rolename SUPERUSER;
user=# ALTER USER rolename WITH PASSWORD 'password';
```

But as documentation sugests

```text
A database superuser bypasses all permission checks. This is a dangerous privilege and should not be used carelessly; it is best to do most of your work as a role that is not a superuser. To create a new database superuser, use CREATE ROLE name SUPERUSER. You must do this as a role that is already a superuser.
```

So this path is not recommended, at least in production environment.

Options:

- Remove permission after running migrations.
- Install POSTGis manually using another allowed user and comment
  - Comment **POSTGis** migrations steps in `internal/mig/mig.go`.

```go
// GetMigrator configured.
func (m *Migrator) addSteps() {
	// Migrations
	// Enable Postgis
	s := &step{}
	s.Config(s.EnablePostgis, s.DropPostgis) // <- comment this line
	m.AddMigration(s) // <- comment this line
```

### Create database

```shell
user=# CREATE DATABASE dbname OWNER rolename;
user=# CREATE DATABASE dbname_test OWNER rolename;
```

Replace dbname by the name of your app database.

### Update run.sh script

Edit `scripts/run.sh`

Replace `baseapp` by chosen app name

```shell
(...)
# Free ports
killall -9 baseapp
(...)
go build -o ./bin/baseapp ./cmd/baseapp.go
./bin/baseapp
```

Config system uses envar prefixes to set app configuration values.
By default this value is `KBS` but you can replace it with any other.

```shell
# Service
export KBS_SVC_NAME="kabestan"
export KBS_SVC_REVISION=$REV
export KBS_(...)
```

Edit other values according to the preferred ones and / or those of your system.

### Edit main

First rename `cmd/baseapp.go` to `cmd/appname.go` where appname is the name you have chosen for your application.

Edit `cmd/appname`

If you change this envvar prefix from "KBS" to, let say, "APP"

```go
  (...)
	cfg := kbs.LoadConfig("kbs") // <- change this
	// cfg := kbs.LoadConfig("app") // <- to something like this
  (...)
```

You can also change, but this value is not used for configuration purposes.

```go
const (
	// Replace by prefered
	appName = "kbs"
)
```

### Run app

```shell
$ make clean-and-run
```

**You should see something like this**

```shell
./scripts/run.sh
2:29AM INF Cookie store key value=hVuOOv4PNBnqTk2o13JsBMOPcPAe4p18
2:29AM INF Reading template path=account/_ctxbar.tmpl
2:29AM INF Reading template path=account/_flash.tmpl
2:29AM INF Reading template path=account/_header.tmpl
2:29AM INF Reading template path=account/_roles.tmpl
2:29AM INF Reading template path=account/roles.tmpl
2:29AM INF Reading template path=layout/base.tmpl
2:29AM INF Reading template path=permission/_ctxbar.tmpl
2:29AM INF Reading template path=permission/_flash.tmpl
2:29AM INF Reading template path=permission/_form.tmpl
2:29AM INF Reading template path=permission/_header.tmpl
2:29AM INF Reading template path=permission/_item.tmpl
2:29AM INF Reading template path=permission/_list.tmpl
2:29AM INF Reading template path=permission/edit.tmpl
2:29AM INF Reading template path=permission/index.tmpl
2:29AM INF Reading template path=permission/initdel.tmpl
2:29AM INF Reading template path=permission/new.tmpl
2:29AM INF Reading template path=permission/show.tmpl
2:29AM INF Reading template path=resource/_ctxbar.tmpl
2:29AM INF Reading template path=resource/_flash.tmpl
2:29AM INF Reading template path=resource/_form.tmpl
2:29AM INF Reading template path=resource/_header.tmpl
2:29AM INF Reading template path=resource/_item.tmpl
2:29AM INF Reading template path=resource/_list.tmpl
2:29AM INF Reading template path=resource/edit.tmpl
2:29AM INF Reading template path=resource/index.tmpl
2:29AM INF Reading template path=resource/initdel.tmpl
2:29AM INF Reading template path=resource/new.tmpl
2:29AM INF Reading template path=resource/show.tmpl
2:29AM INF Reading template path=role/_ctxbar.tmpl
2:29AM INF Reading template path=role/_flash.tmpl
2:29AM INF Reading template path=role/_form.tmpl
2:29AM INF Reading template path=role/_header.tmpl
2:29AM INF Reading template path=role/_item.tmpl
2:29AM INF Reading template path=role/_list.tmpl
2:29AM INF Reading template path=role/edit.tmpl
2:29AM INF Reading template path=role/index.tmpl
2:29AM INF Reading template path=role/initdel.tmpl
2:29AM INF Reading template path=role/new.tmpl
2:29AM INF Reading template path=role/show.tmpl
2:29AM INF Reading template path=user/_ctxbar.tmpl
2:29AM INF Reading template path=user/_flash.tmpl
2:29AM INF Reading template path=user/_form.tmpl
2:29AM INF Reading template path=user/_header.tmpl
2:29AM INF Reading template path=user/_item.tmpl
2:29AM INF Reading template path=user/_list.tmpl
2:29AM INF Reading template path=user/_signin.tmpl
2:29AM INF Reading template path=user/_signup.tmpl
2:29AM INF Reading template path=user/edit.tmpl
2:29AM INF Reading template path=user/index.tmpl
2:29AM INF Reading template path=user/initdel.tmpl
2:29AM INF Reading template path=user/new.tmpl
2:29AM INF Reading template path=user/show.tmpl
2:29AM INF Reading template path=user/signin.tmpl
2:29AM INF Reading template path=user/signup.tmpl
2:29AM INF Parsed template set path=resource/initdel.tmpl
2:29AM INF Parsed template set path=resource/show.tmpl
2:29AM INF Parsed template set path=resource/edit.tmpl
2:29AM INF Parsed template set path=resource/index.tmpl
2:29AM INF Parsed template set path=resource/new.tmpl
2:29AM INF Parsed template set path=permission/show.tmpl
2:29AM INF Parsed template set path=permission/new.tmpl
2:29AM INF Parsed template set path=permission/index.tmpl
2:29AM INF Parsed template set path=permission/edit.tmpl
2:29AM INF Parsed template set path=permission/initdel.tmpl
2:29AM INF Parsed template set path=user/edit.tmpl
2:29AM INF Parsed template set path=user/new.tmpl
2:29AM INF Parsed template set path=user/signup.tmpl
2:29AM INF Parsed template set path=user/index.tmpl
2:29AM INF Parsed template set path=user/initdel.tmpl
2:29AM INF Parsed template set path=user/signin.tmpl
2:29AM INF Parsed template set path=user/show.tmpl
2:29AM INF Parsed template set path=account/roles.tmpl
2:29AM INF Parsed template set path=role/edit.tmpl
2:29AM INF Parsed template set path=role/new.tmpl
2:29AM INF Parsed template set path=role/index.tmpl
2:29AM INF Parsed template set path=role/initdel.tmpl
2:29AM INF Parsed template set path=role/show.tmpl
2:29AM INF Dialing to Postgres host="host=localhost port=5432 user=kabestan password=kabestan dbname=kabestan_dev sslmode=disable"
2:29AM INF Postgres connection established
2:29AM INF New migrator name=migrator
2:29AM INF New seeder name=seeder
2:29AM INF New handler name=ses-mailer
2020/01/25 02:29:26 Migration executed: EnablePostgis
2020/01/25 02:29:27 Migration executed: CreateUsersTable
2020/01/25 02:29:27 Migration executed: CreateAccountsTable
2020/01/25 02:29:27 Migration executed: CreateResourcesTable
2020/01/25 02:29:27 Migration executed: CreateRolesTable
2020/01/25 02:29:27 Migration executed: CreatePermissionsTable
2020/01/25 02:29:27 Migration executed: CreateAccountRolesTable
2020/01/25 02:29:27 Seed step executed: CreateUsersAndAccounts
2020/01/25 02:29:27 Seed step executed: CreateResources
2020/01/25 02:29:27 Seed step executed: CreateRoles
2020/01/25 02:29:27 Seed step executed: CreatePermissions
2020/01/25 02:29:27 Seed step executed: CreateAccountRoles
2:29AM INF Web server initializing port=:8080
```

## Make commands

A brief summary of the most used commands

**make build**

Builds the application

**make run**

Run the application through a shell script that previously sets the environment variables with required values.
In case you need to change some envar, you can edit this script: `scripts/run.sh`.

**make test**

Run tests

**make grc-test**

Run tests with coloured output. [grc](https://github.com/garabik/grc) needs to be available in your system.

**package-resources**

It generates a binary representation for html templates, translations and other resources that allows compiler to embed them within the target file. `clean-and-run` runs this make task as subtask before starting the applicacion.

**build-stg**

Build a staging Docker image of this application and pushes it to Docker Hub.

**build-prod**

Same as `make build-stage` but for production images.

**install-stg**

Deploys app to Googke GKE usando [HELM](https://helm.sh/).
I haven't created helm .yaml files so this command is not functional yet.

**install-prod**

Same as `make install-stage` but for production images.

**TODO**

[GTD](docs/gtd/gtd.md)
