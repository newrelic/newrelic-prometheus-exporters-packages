SHELL := /bin/bash
NRI_GENERATOR_PATH="$(PWD)/nri-config-generator"
INTEGRATIONS="ravendb github githubactions"


clean:
	rm -rf dist

build-all:
	for integration in "$(INTEGRATIONS)"; do \
		make 	build-$${integration}; \
	done;

build-%:
	sh scripts/clean.sh $(PWD) $*
	sh exporters/$*/build.sh $(PWD)
	sh scripts/build_generator.sh $(PWD) $*
	sh scripts/copy_resources.sh $(PWD) $*



package-%:
	sh scripts/package.sh $(PWD) $*

all:
	for integration in "$(INTEGRATIONS)"; do \
		make	build-$${integration}; \
		make 	package-$${integration}; \
	done;

run:
	sh scripts/run.sh $(PWD)
	docker-compose -f tests/docker-compose.yml up

include $(CURDIR)/nri-config-generator/build/ci.mk