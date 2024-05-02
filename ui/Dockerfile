FROM node:20.12.2-alpine3.19 AS base

RUN apk add curl

# Entrypoint and entrypoint utils
RUN apk add --no-cache python3 py3-pip py3-watchdog curl
RUN mkdir -p /utils
COPY ./utils /utils
COPY ./ui/entrypoint /

# Necessary for devscripts/openapi-spec
RUN apk add openjdk11-jre
RUN npm install -g @openapitools/openapi-generator-cli
RUN openapi-generator-cli version-manager set 7.5.0

ENV NEXT_TELEMETRY_DISABLED 1
ENV PORT 8080

##
##
##
FROM base AS deps

WORKDIR /app

RUN apk add --no-cache libc6-compat

COPY ./ui/package.json ./ui/pnpm-lock.yaml* ./
RUN corepack enable pnpm && \
    pnpm i --frozen-lockfile

##
##
##
FROM base as development

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY ./ui .

ENTRYPOINT ["/entrypoint"]

##
##
##
FROM base AS builder

COPY --from=deps /app/node_modules ./node_modules
COPY ./ui .

RUN corepack enable pnpm && \
    pnpm run build

##
##
##
FROM base as production

WORKDIR /app

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public

RUN mkdir .next
RUN chown nextjs:nodejs .next

COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

ENTRYPOINT ["/entrypoint"]
