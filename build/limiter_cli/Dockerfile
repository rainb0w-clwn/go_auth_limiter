FROM golang:1.25-alpine3.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o cli ./cmd/limiter_cli

COPY ./configs/config.yml /configs/config.yml

ENTRYPOINT ["/app/cli"]
