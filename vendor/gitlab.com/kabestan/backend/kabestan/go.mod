module gitlab.com/kabestan/backend/kabestan

go 1.13

require (
	github.com/aws/aws-sdk-go v1.28.7
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/go-chi/chi v4.0.3+incompatible
	github.com/gorilla/csrf v1.6.2
	github.com/gorilla/schema v1.1.0
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.3.0
	github.com/markbates/pkger v0.14.0
	github.com/nicksnyder/go-i18n/v2 v2.0.3
	github.com/rs/zerolog v1.17.2
	github.com/satori/go.uuid v1.2.0
	gitlab.com/kabestan/backend/kabestan/db v0.0.0-20200111024303-1eea754ea2f3
	gitlab.com/kabestan/backend/kabestan/db/pg v0.0.0-20200111024303-1eea754ea2f3
	golang.org/x/text v0.3.2
	google.golang.org/appengine v1.6.5 // indirect
)
