SHELL := /bin/bash
NRI_GENERATOR_PATH="$(PWD)/nri-config-generator"

NEWRELIC_E2E ?= go run github.com/newrelic/newrelic-integration-e2e-action@latest
GOOS ?= "linux"

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
	bash exporters/$*/build-${GOOS}.sh $(PWD) && \
	bash scripts/build_generator.sh $(PWD) $* ${GOOS};

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
	GOOS=linux make build-$*
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
