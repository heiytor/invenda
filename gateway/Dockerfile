FROM nginx:1.26.0-alpine3.19 AS base

RUN apk add --no-cache python3 py3-pip

RUN ["rm", "/etc/nginx/conf.d/default.conf"]

RUN mkdir -p /etc/nginx/default.d
RUN mkdir -p /var/run/nginx/

WORKDIR /etc/nginx
COPY ./nginx.conf /etc/nginx/

COPY ./entrypoint /
ENTRYPOINT ["/entrypoint"]
CMD ["/usr/sbin/nginx", "-g", "'daemon off;'"]

FROM base as development

# ...

FROM base as production

# ...
