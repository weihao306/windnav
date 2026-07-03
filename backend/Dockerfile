FROM golang:1.26-alpine AS backend-builder

WORKDIR /app/backend
RUN apk add --no-cache build-base
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=1 go build -o /out/windnav ./cmd/server

FROM node:24-alpine AS frontend-builder

WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM alpine:3.23

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=backend-builder /out/windnav /app/windnav
COPY --from=frontend-builder /app/frontend/dist /app/public
RUN mkdir -p /app/data

ENV HTTP_ADDR=:8080
ENV DB_DRIVER=sqlite
ENV DB_SQLITE_PATH=/app/data/windnav.db
ENV STATIC_DIR=/app/public

EXPOSE 8080
CMD ["/app/windnav"]
