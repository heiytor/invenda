FROM node:20.12.2-alpine3.19 AS base

RUN apk add curl

# Entrypoint and entrypoint utils
RUN apk add --no-cache python3 py3-pip py3-watchdog curl
RUN mkdir -p /utils
COPY ./utils /utils
COPY openapi/entrypoint /

RUN npm install -g @redocly/cli@1.0.0-beta.100

WORKDIR /openapi
COPY openapi/package*.json ./

ENTRYPOINT ["/entrypoint"]

FROM base as development

# ...

FROM base as production

# ...
