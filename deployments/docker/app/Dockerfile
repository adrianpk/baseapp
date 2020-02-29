# Base build image
FROM golang:1.12-alpine AS build

# Install Git
RUN apk update && apk add git

# Set workdir
WORKDIR /go/src/gitlab.com/kabestan/repo/baseapp

# Force the go compiler to use modules
ENV GO111MODULE=on

# Copy project content
RUN echo "Copying resources..."
COPY . .
RUN echo "Copied."

# Compile the project
ENV CGO_ENABLED=0 GOOS=linux
RUN echo "Compiling..."
RUN go build -mod=vendor -o /go/bin/baseapp cmd/baseapp.go
RUN echo "Compile output in: /go/bin/baseapp"

# Coying only the results without the artifacts to fresh Alpine image.
FROM alpine
WORKDIR /srv
COPY --from=build /go/bin/baseapp ./baseapp

# Entrypoint
ENTRYPOINT ./baseapp
