FROM golang:1.24.4-alpine as builder

WORKDIR /app


COPY . .

RUN go mod download

RUN go build -v -o /app/uwud .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/uwud /app/uwud
COPY --from=builder /app/views  /app/views
COPY --from=builder /app/static /app/static

RUN chmod +x uwud

ENV GIN_MODE=release

CMD ["./uwud"]
