FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/commlayers-api ./cmd/server

FROM alpine:3.20

RUN adduser -D -H -u 10001 appuser

COPY --from=builder /bin/commlayers-api /usr/local/bin/commlayers-api

USER appuser

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/commlayers-api"]
