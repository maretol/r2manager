FROM node:22 AS base
FROM base AS deps

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci


FROM base AS builder

ARG BASE_PATH=""

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

ENV NEXT_TELEMETRY_DISABLED=1

RUN echo "BACKEND_URL=http://localhost:8080" >> .env
RUN echo "HOST=http://localhost:3000" >> .env
RUN echo "BASE_PATH=${BASE_PATH}" >> .env
RUN echo "NEXT_PUBLIC_APP_NAME='R2 Manager'" >> .env

RUN npm run build


FROM node:22-slim AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# default env vars
ENV BACKEND_URL="http://localhost:8080"
ENV HOSTNAME="0.0.0.0"
ENV PROTOCOL="http"
ENV PORT="3000"

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENTRYPOINT ["node", "server.js"]
