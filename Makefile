SHELL := /bin/bash
NRI_GENERATOR_PATH="$(PWD)/nri-config-generator"

GOOS ?=
GOARCH ?=

NEWRELIC_E2E ?= go run github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e@latest

clean:
	rm -rf dist

build-all:
	@cd exporters; \
	for name in $$(ls -d *) ; do \
		cd $(PWD); \
		make build-$${name}; \
	done

build-%: GOOS := $(if $(GOOS),$(GOOS),"linux")
build-%: GOARCH := $(if $(GOARCH),$(GOARCH),"arm64")
build-%:
	@echo "[ build-$* ]: Building exporter..."
	source scripts/common_functions.sh; \
	EXPORTER_PATH=exporters/$*/exporter.yml; \
	loadVariables; \
	bash exporters/$*/build.sh $(PWD) ${GOOS} ${GOARCH} && \
	bash scripts/build_generator.sh $(PWD) $*;

fetch-resources-%:
	@echo "[ fetch-resources-$* ]: Fetching external resources..."
	source scripts/common_functions.sh; \
	EXPORTER_PATH=exporters/$*/exporter.yml; \
	loadVariables; \
	bash scripts/create_folder_structure.sh $(PWD) $* && \
	bash scripts/fetch_external_files.sh $(PWD) $*;

package-%:
	@echo "[ package-$* ]: Packaging exporter..."
	source scripts/common_functions.sh; \
	EXPORTER_PATH=exporters/$*/exporter.yml; \
	loadVariables; \
	bash scripts/create_folder_structure.sh $(PWD) $* && \
	bash scripts/copy_resources.sh $(PWD) $* && \
	bash scripts/package.sh $(PWD) $*

test-e2e-%: build-%
	@echo "[ test-e2e-%$* ]: Running e2e test..."
	$(NEWRELIC_E2E) --commit_sha=test-string --retry_attempts=5 --retry_seconds=60 \
         --account_id=$(ACCOUNT_ID) --api_key=$(API_KEY) --license_key=$(LICENSE_KEY) \
         --spec_path=$(PWD)/exporters/$*/e2e/e2e_spec.yml --verbose_mode

all:
	@cd exporters; \
	for name in $$(ls -d *) ; do \
		cd $(PWD); \
		make	build-$${name} && \
		make	fetch-resources-$${name} && \
		make 	package-$${name}; \
	done

run:
	bash scripts/run.sh $(PWD)
	docker-compose -f tests/docker-compose.yml up

include $(CURDIR)/nri-config-generator/build/ci.mk
