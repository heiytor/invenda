FROM nginx:1.26.0-alpine3.19 AS base

RUN apk add curl

RUN ["rm", "/etc/nginx/conf.d/default.conf"]

# Entrypoint and entrypoint utils
RUN apk add --no-cache python3 py3-pip py3-watchdog
RUN mkdir -p /utils
COPY ./utils /utils
COPY ./gateway/entrypoint /

COPY --from=hairyhenderson/gomplate:v2.5.0-slim /gomplate /bin/gomplate

RUN mkdir -p /etc/nginx/default.d
RUN mkdir -p /var/run/nginx/

WORKDIR /etc/nginx
COPY ./gateway/nginx.conf /etc/nginx/spec.conf

ENTRYPOINT ["/entrypoint"]

FROM base as development

# ...

FROM base as production

# ...
