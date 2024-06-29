FROM golang:1.22-alpine AS builder

RUN apk update && apk add ca-certificates
RUN update-ca-certificates

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o main ./cmd/bot/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/main /main
COPY --from=builder /build/config/config-local.yml /config/config-local.yml

ENTRYPOINT ["/main"]