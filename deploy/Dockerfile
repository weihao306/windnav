FROM alpine:3.24 AS backend-builder
WORKDIR /app/backend
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk add --no-cache build-base go
COPY backend/ ./
RUN CGO_ENABLED=1 go build -mod=vendor -o /out/windnav ./cmd/server

FROM alpine:3.24 AS frontend-builder
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk add --no-cache build-base nodejs npm
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM alpine:3.24 AS runner
WORKDIR /app
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk add --no-cache ca-certificates tzdata
COPY --from=backend-builder /out/windnav /app/windnav
COPY --from=frontend-builder /app/frontend/dist /app/public
RUN mkdir -p /app/data

ENV HTTP_ADDR=:8080
ENV DB_DRIVER=sqlite
ENV DB_SQLITE_PATH=/app/data/windnav.db
ENV STATIC_DIR=/app/public

EXPOSE 8080
CMD ["/app/windnav"]
