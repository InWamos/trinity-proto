FROM golang:1.25.7-trixie AS builder

ARG CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o trinity ./cmd/main.go

FROM scratch
COPY --from=builder /app/trinity /trinity
ENTRYPOINT ["/trinity"]