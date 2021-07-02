TARGET          := target
GOFLAGS          = -mod=readonly
GOLANGCI_LINT	 = github.com/golangci/golangci-lint/cmd/golangci-lint

all: build

build: clean validate compile test

clean:
	@echo "=== nri-config-generator === [ clean ]: removing binaries and coverage file..."
	@rm -rfv bin coverage.xml $(TARGET)

validate:
	@printf "=== nri-config-generator === [ validate ]: running golangci-lint & semgrep... "
	@go run  $(GOFLAGS) $(GOLANGCI_LINT) run --verbose
	@[ -f .semgrep.yml ] && semgrep_config=".semgrep.yml" || semgrep_config="p/golang" ; \
	docker run --rm -v "${PWD}:/src:ro" --workdir / returntocorp/semgrep -c "$$semgrep_config"

test:
	@echo "=== nri-config-generator === [ test ]: running unit tests..."
	@go test -race ./... -count=1

integration-test:
	@echo "=== nri-config-generator === [ integration-test ]: running integration tests..."
	@go test -v -tags=integration ./tests/.

PACKAGE_NAME?=
VERSION?=dev
GIT_COMMIT?=$(shell git rev-parse HEAD)
BUILD_DATE?=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
BIN_PATH?=bin/$(PACKAGE_NAME)

compile:
ifeq ($(PACKAGE_NAME),)
	$(error Missing required argument PACKAGE_NAME.)
endif
	go build \
	-ldflags "-X main.integration=$(PACKAGE_NAME) -X main.integrationVersion=$(VERSION) -X main.gitCommit=$(GIT_COMMIT) -X main.buildDate=$(BUILD_DATE)" \
	-o $(BIN_PATH) github.com/newrelic/nri-config-generator

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
PARENT_DIR:=$(shell cd $(CURDIR)/..; pwd)

fetch-definitions:
ifeq ($(PACKAGE_NAME),)
	$(error Missing required argument PACKAGE_NAME.)
endif
	@echo "[ build-$* ]: Building exporter...$(PACKAGE_NAME)"
	cd $(PARENT_DIR) ; \
		make build-$(PACKAGE_NAME)