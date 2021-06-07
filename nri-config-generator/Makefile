INTEGRATION      = nri-config-generator
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


nri-config-generator/snyk-test:
	@docker run --rm -t \
			--name "nri-config-generator-snyk-test" \
			-v $(CURDIR):/go/src/github.com/newrelic/nri-config-generator \
			-w /go/src/github.com/newrelic/nri-config-generator \
			-e SNYK_TOKEN \
			-e GO111MODULE=auto \
			snyk/snyk:golang snyk test --severity-threshold=high


INTEGRATION_NAME=
compile:
ifeq (nri-config-generator,)
	$(error Missing required argument INTEGRATION_NAME.)
endif
	go build -ldflags "-X main.integration=$(INTEGRATION_NAME) -o bin/nri-$(INTEGRATION_NAME) github.com/newrelic/nri-config-generator"
	echo "created executable nri-$(INTEGRATION_NAME)"