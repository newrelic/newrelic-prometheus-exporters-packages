SHELL := /bin/bash
NRI_GENERATOR_PATH="$(PWD)/nri-config-generator"

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