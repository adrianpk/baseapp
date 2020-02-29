# Base build image
FROM golang:1.12-alpine AS build

# Install git
RUN set -ex; \
    apk update; \
    apk add --no-cache git

# Set workdir
WORKDIR /go/src/gitlab.com/kabestan/repo/baseapp

# Run tests
# CMD CGO_ENABLED=0 go test ./...
CMD CGO_ENABLED=0 go test -v -run TestCreateUser -count=1 -timeout=5s  ./pkg/web/user_test.go ./pkg/web/base_test.go
