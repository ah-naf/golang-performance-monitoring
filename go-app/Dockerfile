FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -o server main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

ENV GIN_MODE=release

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

ENTRYPOINT ["./server"]