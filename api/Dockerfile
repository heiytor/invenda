# stage used by targets as a base
FROM golang:1.21.6-alpine3.19 AS base

RUN apk add curl

WORKDIR $GOPATH/src/github.com/heiytor/invenda/api

COPY api/go.mod api/go.sum .
RUN go mod download

# development target
FROM base AS dev

# refresh is used to refreshs the server when any file is updated
RUN go install github.com/markbates/refresh@v1.11.1

WORKDIR $GOPATH/src/github.com/heiytor/invenda/api
COPY api . 

COPY api/entrypoint.sh /entrypoint.sh

WORKDIR $GOPATH/src/github.com/heiytor/invenda/api/cmd
ENTRYPOINT /entrypoint.sh

# production target
FROM base AS prod
# TODO

