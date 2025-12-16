# Собираем в гошке
FROM golang:1.25 AS build

ENV BIN_FILE /opt/limiter/limiter-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/limiter/*

# На выходе тонкий образ
FROM alpine:3.23

RUN apk add --no-cache curl

LABEL ORGANIZATION="LittleBig"
LABEL SERVICE="limiter"
LABEL MAINTAINERS="mikhail.anikin15@gmail.com"

ENV BIN_FILE /opt/limiter/limiter-app
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/limiter/config.yml
COPY ./configs/config.yml ${CONFIG_FILE}

COPY migrations /app/migrations

CMD ${BIN_FILE} -config ${CONFIG_FILE}
