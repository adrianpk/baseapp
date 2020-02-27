module gitlab.com/kabestan/repo/baseapp

go 1.13

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/go-chi/chi v4.0.3+incompatible
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.3.0
	github.com/markbates/pkger v0.14.0
	github.com/satori/go.uuid v1.2.0
	github.com/siddontang/go v0.0.0-20180604090527-bdc77568d726
	gitlab.com/kabestan/backend/kabestan v0.0.0
	gitlab.com/kabestan/backend/kabestan/db v0.0.0
	golang.org/x/crypto v0.0.0-20200117160349-530e935923ad
	golang.org/x/text v0.3.2
	honnef.co/go/tools v0.0.1-2020.1.3
)

replace gitlab.com/kabestan/backend/kabestan => ../../backend/kabestan

replace gitlab.com/kabestan/backend/kabestan/db => ../../backend/kabestan/db
