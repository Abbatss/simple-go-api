APP_NAME?=$(shell basename `pwd`)
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
REGISTRY?=test
IMAGE=$(REGISTRY)/$(APP_NAME)
GO_MAIN_SRC?=cmd/server/main.go

TAG_NAME?=$(shell git describe --tags 2> /dev/null | echo SNAPSHOT)
SHORT_SHA?=$(shell git rev-parse --short HEAD)
VERSION?=$(TAG_NAME)-$(SHORT_SHA)


GO_LINT?=golangci-lint
GOCMD?=CGO_ENABLED=0 go
GOCMD_TEST?=GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn APP_ENV=test CGO_ENABLED=0 go

.PHONY: lint
lint: lint_code

.PHONY: lint_code
lint_code:
	$(GO_LINT) -c build/golangci.yaml run

.PHONY: test
test:
	@make up
	@make test_integration
	@make down

.PHONY: migrate
migrate:
		sleep 4;migrate -database postgres://user:password@localhost:5432/test?sslmode=disable -path migrations up

.PHONY: up
up:
	@docker-compose -f build/docker-compose.yaml up -d
	@make migrate
.PHONY: down
down:
	@docker-compose -f build/docker-compose.yaml down -v

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: test_integration
test_integration:
	$(GOCMD_TEST) test ./... -mod=vendor -count=1 -tags=integration

.PHONY: test_unit
test_unit:
	$(GOCMD_TEST) test ./... -mod=vendor -count=1

.PHONY: build
build:
	$(GOCMD) build -mod vendor -ldflags "-X main.serviceVersion=$(VERSION)" -o go-app $(GO_MAIN_SRC)

.PHONY: vendor
vendor:
	$(GOCMD) mod vendor

.PHONY: mod
mod:
	$(GOCMD) get -u
	$(GOCMD) mod tidy
	make vendor

.PHONY: image
image:
	docker build -f build/Dockerfile --build-arg=VERSION=$(VERSION) -t $(IMAGE):$(VERSION) .

.PHONY: clean
clean:
	@$(GOCMD) clean

