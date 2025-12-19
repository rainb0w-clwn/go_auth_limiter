GOLANGCI_LINT_VERSION=v2.7.2

BIN := "./bin/limiter"
CONFIG_PATH := "./configs/config.yml"
COMPOSE_FILE := "deployments/docker-compose.yml"
INTEGRATION_COMPOSE_FILE := "deployments/docker-compose.integration.yml"

generate:
	buf generate

tools:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/meshapi/grpc-api-gateway/codegen/cmd/protoc-gen-openapiv3@latest
	go install github.com/meshapi/grpc-api-gateway/codegen/cmd/protoc-gen-grpc-api-gateway@latest
	buf dep update

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/limiter

start: build
	$(BIN) -config $(CONFIG_PATH)

run:
	docker compose -f $(COMPOSE_FILE) up --build

run-cli:
	docker compose -f $(COMPOSE_FILE) run --rm cli $(ARGS)

down:
	docker compose -f $(COMPOSE_FILE) down

test:
	CGO_ENABLED=1 go test --count=1 -race ./internal/...

integration-test:
	docker compose -f ${INTEGRATION_COMPOSE_FILE} down --volumes
	docker compose -f $(INTEGRATION_COMPOSE_FILE) up --build \
		--abort-on-container-exit \
		--exit-code-from test
	docker compose -f ${INTEGRATION_COMPOSE_FILE} down --rmi local

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}

lint: install-lint-deps
	golangci-lint run ./...

clean:
	rm -f **/*.pb.go
	rm -f **/*.gw.go
	rm -f **/*.openapi.yaml

.PHONY: generate tools build run up down test integration-tests install-lint-deps lint clean
