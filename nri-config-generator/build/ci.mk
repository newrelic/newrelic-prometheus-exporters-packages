INTEGRATION = config-generator
BUILDER_TAG ?= nri-$(INTEGRATION)-builder

SRC_DIR=$(CURDIR)/nri-config-generator

.PHONY : ci/deps
ci/deps:
	@docker build -t $(BUILDER_TAG) -f $(SRC_DIR)/build/Dockerfile $(CURDIR)

ci/snyk-test:
	@docker run --rm -t \
		--name "nri-$(INTEGRATION)-snyk-test" \
		-v $(SRC_DIR):/go/src/github.com/newrelic/nri-config-generator \
		-w /go/src/github.com/newrelic/nri-config-generator \
		-e SNYK_TOKEN \
		-e GO111MODULE=auto \
		snyk/snyk:golang snyk test --severity-threshold=high

.PHONY : ci/validate
ci/validate: ci/deps
	@docker run --rm -t \
			--name "nri-$(INTEGRATION)-validate" \
			-v $(SRC_DIR):/go/src/github.com/newrelic/nri-$(INTEGRATION) \
			-w /go/src/github.com/newrelic/nri-$(INTEGRATION) \
			$(BUILDER_TAG) make validate

.PHONY : ci/test
ci/test: ci/deps
	@docker run --rm -t \
		--name "nri-$(INTEGRATION)-test" \
		-v $(SRC_DIR):/go/src/github.com/newrelic/nri-$(INTEGRATION) \
		-w /go/src/github.com/newrelic/nri-$(INTEGRATION) \
		$(BUILDER_TAG) make test
