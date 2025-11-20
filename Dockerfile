FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o pq-app ./cmd/question_answer


FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/pq-app /app/pq-app
COPY --from=builder /app/config /app/config
COPY --from=builder /app/internal/infrastructure/storage/postgres/migrations \
    /app/internal/infrastructure/storage/postgres/migrations

EXPOSE 8082
CMD ["./pq-app"]