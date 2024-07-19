FROM golang:1.20-buster AS builder

RUN apt-get update && apt-get install -y ca-certificates openssl
RUN update-ca-certificates

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o main .

FROM debian:buster-slim
WORKDIR /

COPY --from=builder /app/main /
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

ENTRYPOINT ["/main"]
CMD ["service"]