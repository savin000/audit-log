FROM golang:1.24.3 AS builder

ARG CGO_ENABLED=0
WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY . .

RUN go build -o audit-log ./cmd/main.go

FROM scratch
COPY --from=builder /app/audit-log /audit-log
ENTRYPOINT ["/audit-log"]