FROM golang:1.26-alpine AS backend-builder
WORKDIR /app/backend
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk add --no-cache build-base
COPY backend/ ./
RUN CGO_ENABLED=1 go build -mod=vendor -o /out/windnav ./cmd/server

FROM node:24-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM nginx:1.29-alpine AS runner
WORKDIR /app
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk add --no-cache ca-certificates tzdata
COPY --from=backend-builder /out/windnav /app/windnav
COPY --from=frontend-builder /app/frontend/dist /usr/share/nginx/html
COPY deploy/nginx.conf /etc/nginx/conf.d/default.conf
COPY deploy/docker-entrypoint.sh /docker-entrypoint.sh
RUN mkdir -p /app/data

RUN chmod +x /docker-entrypoint.sh

ENV HTTP_ADDR=:8080
ENV DB_DRIVER=sqlite
ENV DB_SQLITE_PATH=/app/data/windnav.db

EXPOSE 80
ENTRYPOINT ["/docker-entrypoint.sh"]
