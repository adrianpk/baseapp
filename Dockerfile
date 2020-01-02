# Base build image
FROM golang:1.11-alpine AS build

# Install Git
RUN apk update && apk add git

# Set workdir
WORKDIR /go/src/gitlab.com/kabestan/backend/kabestan

# Force the go compiler to use modules
ENV GO111MODULE=on

# Copy project content
RUN echo "Copying resources..."
COPY . .
RUN echo "Copied."

# Compile the project
ENV CGO_ENABLED=0 GOOS=linux
RUN echo "Compiling..."
RUN go build -mod=vendor -o /go/bin/kabestan cmd/kabestan.go
RUN echo "Compile output in: /go/bin/kabestan"

# Coying only the results without the artifacts to fresh Alpine image.
FROM alpine
WORKDIR /srv
COPY --from=build /go/bin/kabestan ./kabestan

# Entrypoint
ENTRYPOINT ./kabestan
