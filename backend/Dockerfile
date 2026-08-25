# =========================
# Stage 1 — Build
# =========================

FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o mindk-api ./cmd/api


# =========================
# Stage 2 — Runtime
# =========================

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/mindk-api .

EXPOSE 8080

CMD ["./mindk-api"]