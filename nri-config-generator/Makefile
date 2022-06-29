TARGET          := target

all: build

build: clean compile test

clean:
	@echo "=== nri-config-generator === [ clean ]: removing binaries and coverage file..."
	@rm -rfv bin coverage.xml $(TARGET)

test:
	@echo "=== nri-config-generator === [ test ]: running unit tests..."
	@go test -race ./... -count=1

integration-test:
	@echo "=== nri-config-generator === [ integration-test ]: running integration tests..."
	@go test -v -tags=integration ./integration-tests/.

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
