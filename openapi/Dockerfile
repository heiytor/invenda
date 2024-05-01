FROM node:20.12.2-alpine3.19 AS base

RUN apk add inotify-tools curl
RUN npm install -g @redocly/cli@1.0.0-beta.100

WORKDIR /openapi
COPY package*.json ./

ENTRYPOINT ["/openapi/entrypoint.sh"]
