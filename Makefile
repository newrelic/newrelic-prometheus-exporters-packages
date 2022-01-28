SHELL := /bin/bash
NRI_GENERATOR_PATH="$(PWD)/nri-config-generator"

ifeq (, $(shell which newrelic-integration-e2e))
NEWRELIC_E2E ?= go run github.com/newrelic/newrelic-integration-e2e-action/newrelic-integration-e2e/cmd@v1
else
NEWRELIC_E2E ?= newrelic-integration-e2e
endif

clean:
	rm -rf dist

build-all:
	@cd exporters; \
	for name in $$(ls -d *) ; do \
		cd $(PWD); \
		make build-$${name}; \
	done

build-%:
	@echo "[ build-$* ]: Building exporter..."
	source scripts/common_functions.sh; \
	EXPORTER_PATH=exporters/$*/exporter.yml; \
	loadVariables; \
	bash scripts/fetch_synthesis_definition.sh $(PWD) && \
	bash exporters/$*/build.sh $(PWD) && \
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

test-e2e-%: 
	@echo "[ test-e2e-%$* ]: Running e2e test..."
	$(NEWRELIC_E2E) --commit_sha=test-string --retry_attempts=5 --retry_seconds=60 \
         --account_id=$(ACCOUNT_ID) --api_key=$(API_KEY) --license_key=$(LICENSE_KEY) \
         --spec_path=$(PWD)/exporters/$*/e2e/e2e_spec.yml --verbose_mode=true

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