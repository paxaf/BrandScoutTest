FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN go mod download

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# RUN go test ./... -race -cover -v
RUN go build -o quotes_app ./cmd

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/quotes_app quotes_app

ENTRYPOINT ["./quotes_app"]

