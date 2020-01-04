#!/bin/sh
# Build
# ./scripts/build.sh

# Free ports
killall -9 baseapp

# Set environment variables
REV=$(eval git rev-parse HEAD)
# Service
export KBS_SVC_NAME="kabestan"
export KBS_SVC_REVISION=$REV
export KBS_SVC_PINGPORT=8090
# Servers
export KBS_WEB_SERVER_PORT=8080
export KBS_JSONREST_SERVER_PORT=8081
export KBS_WEB_COOKIESTORE_KEY="hVuOOv4PNBnqTk2o13JsBMOPcPAe4p18"
export KBS_WEB_SECCOOKIE_HASH="hVuOOv4PNBnqTk2o13JsBMOPcPAe4p18"
export KBS_WEB_SECCOOKIE_BLOCK="hVuOOv4PNBnqTk2o"
export KBS_SITE_URL="localhost:8080"
# Postgres
export KBS_PG_SCHEMA="public"
export KBS_PG_DATABASE="kabestan_dev"
export KBS_PG_HOST="localhost"
export KBS_PG_PORT="5432"
export KBS_PG_USER="kabestan"
export KBS_PG_PASSWORD="kabestan"
export KBS_PG_BACKOFF_MAXTRIES="3"
# Confirmation
## users/{slug}/{token}/confirm
export KBS_USER_CONFIRMATION_PATH="users/%s/%s/confirm"
export KBS_USER_CONFIRMATION_SEND="false"
export KBS_USER_CONFIRMATION_DEBUG="true"
# Amazon SES MAiler
  # These are sample not usable keys
export AWS_ACCESS_KEY_ID=EIIAHI5FF3A2OG3MJEX5
export AWS_SECRET_KEY=8BiWmd5Hdgmk2rR4pwG332bHwvLGiJOoxLLtDy12

# Switches
export KBS_APP_USERNAME_UPDATABLE=false

go build -o ./bin/baseapp ./cmd/baseapp.go
./bin/baseapp
# go run -race main.go
