# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /jsonstat ./cmd/jsonstat/


# Run stage
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /jsonstat .

ENTRYPOINT ["./jsonstat"]