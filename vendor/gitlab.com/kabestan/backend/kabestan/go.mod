module gitlab.com/kabestan/backend/kabestan

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/aws/aws-sdk-go v1.26.8
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/gorilla/csrf v1.6.2
	github.com/gorilla/schema v1.1.0
	github.com/gorilla/sessions v1.2.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.3.0
	github.com/markbates/pkger v0.13.0
	github.com/nicksnyder/go-i18n/v2 v2.0.3
	github.com/rs/zerolog v1.17.2
	github.com/satori/go.uuid v1.2.0
	gitlab.com/kabestan/backend/kabestan/db v0.0.0-00010101000000-000000000000
	gitlab.com/kabestan/backend/kabestan/db/pg v0.0.0-00010101000000-000000000000
	golang.org/x/text v0.3.2
	google.golang.org/appengine v1.6.5 // indirect
)

replace gitlab.com/kabestan/backend/kabestan/db => ./db

replace gitlab.com/kabestan/backend/kabestan/db/pg => ./db/pg
