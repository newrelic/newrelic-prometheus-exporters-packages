SHELL := /bin/bash
NRI_GENERATOR_PATH="$(PWD)/nri-config-generator"
GORELEASER_VERSION := v1.13.1
GORELEASER_BIN ?= bin/goreleaser

NEWRELIC_E2E ?= go run github.com/newrelic/newrelic-integration-e2e-action@latest
GOOS ?= "linux"

bin:
	@mkdir -p "./bin"

bin/goreleaser: bin
	@echo "[$(GORELEASER_BIN)] Installing goreleaser $(GORELEASER_VERSION)"
	@wget -qO /tmp/goreleaser.tar.gz https://github.com/goreleaser/goreleaser/releases/download/$(GORELEASER_VERSION)/goreleaser_$(OS_DOWNLOAD)_x86_64.tar.gz
	@tar -xf  /tmp/goreleaser.tar.gz -C bin/
	@rm -f /tmp/goreleaser.tar.gz
	@echo "[$(GORELEASER_BIN)] goreleaser downloaded"

clean:
	@rm -rf dist

build-%: clean bin/goreleaser
	@echo "[ build-$* ]: Building exporter..."
	bash ./scripts/build.sh $(PWD) $* $(GOOS)

create-publish-schema-%:
	@echo "[ publish-schema ]: Creating publish schema..."
	bash ./scripts/create_publish_schema.sh $(PWD) $*

package-%: clean bin/goreleaser build-%
	@echo "[ package-$* ]: Packaging exporter..."
	bash ./scripts/package.sh $(PWD) $* $(GOOS)

test-e2e-%:
	@echo "[ test-e2e-%$* ]: Running e2e test..."
	@GOOS=linux make build-$*
	$(NEWRELIC_E2E) --commit_sha=test-string --retry_attempts=5 --retry_seconds=60 \
         --account_id=$(ACCOUNT_ID) --api_key=$(API_KEY) --license_key=$(LICENSE_KEY) \
         --spec_path=$(PWD)/exporters/$*/e2e/e2e_spec.yml --verbose_mode=true

OS := $(shell uname -s)
ifeq ($(OS), Darwin)
	OS_DOWNLOAD := "darwin"
	TAR := gtar
else
	OS_DOWNLOAD := "linux"
endif
