FROM golang:1.24.4-alpine AS builder

WORKDIR /app


COPY . .

RUN go mod download

RUN go build -v -o /app/uwud .

FROM scratch

WORKDIR /app

COPY --from=builder /app/uwud /app/uwud
COPY --from=builder /app/views  /app/views
COPY --from=builder /app/static /app/static


ENV GIN_MODE=release

CMD ["./uwud"]
