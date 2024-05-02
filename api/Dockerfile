# stage used by targets as a base
FROM golang:1.22.2-alpine3.19 AS base

RUN apk add curl

# Entrypoint and entrypoint utils
RUN apk add --no-cache python3 py3-pip py3-watchdog curl
RUN mkdir -p /utils
COPY ./utils /utils
COPY api/entrypoint /

WORKDIR $GOPATH/src/github.com/heiytor/invenda/api

# Secret pairs
RUN mkdir -p ./../secrets
COPY secrets/api-private.pem ./../secrets
COPY secrets/api-public.pem ./../secrets

# Download go packages
COPY api/go.mod api/go.sum .
RUN go mod download

ENTRYPOINT ["/entrypoint"]

FROM base AS development

# Download air and sync files to live reload
WORKDIR $GOPATH/src/github.com/heiytor/invenda/api/
RUN go install github.com/cosmtrek/air@latest
COPY api . 

FROM base AS production

# ...
