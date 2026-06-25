FROM golang:1.22-alpine AS builder

# CGO requires gcc and musl-dev on alpine
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=1 is required for go-sqlite3 (uses cgo bindings)
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o student-api ./cmd/student-api

# ---- final image ----
FROM alpine:3.20

# ca-certificates for HTTPS calls, sqlite-libs for the sqlite3 shared lib
RUN apk add --no-cache ca-certificates sqlite-libs

WORKDIR /app

RUN mkdir -p /app/storage

COPY --from=builder /app/student-api /app/student-api
COPY config ./config

ENV CONFIG_PATH=/app/config/local.yaml

EXPOSE 3000

CMD ["/app/student-api"]
