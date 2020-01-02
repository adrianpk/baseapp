module gitlab.com/kabestan/repo/baseapp

go 1.13

require (
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.3.0
	github.com/markbates/pkger v0.13.0
	github.com/nicksnyder/go-i18n/v2 v2.0.3
	github.com/satori/go.uuid v1.2.0
	gitlab.com/kabestan/backend/kabestan v0.0.0-20191228021211-f4dcb668bc31
	gitlab.com/kabestan/backend/kabestan/db v0.0.0
	gitlab.com/kabestan/backend/kabestan/db/pg v0.0.0 // indirect
	golang.org/x/crypto v0.0.0-20191227163750-53104e6ec876
	golang.org/x/text v0.3.2
)

replace gitlab.com/kabestan/backend/kabestan => ../../backend/kabestan

replace gitlab.com/kabestan/backend/kabestan/db => ../../backend/kabestan/db

replace gitlab.com/kabestan/backend/kabestan/db/pg => ../../backend/kabestan/db/pg
