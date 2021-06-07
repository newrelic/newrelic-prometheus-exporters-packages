INTEGRATION = config-generator
BUILDER_TAG ?= nri-$(INTEGRATION)-builder

.PHONY : ci/deps
ci/deps:
	@docker build -t $(BUILDER_TAG) -f $(CURDIR)/build/Dockerfile $(CURDIR)

ci/snyk-test:
	@docker run --rm -t \
		--name "nri-$(INTEGRATION)-snyk-test" \
		-v $(CURDIR):/go/src/github.com/newrelic/nri-config-generator \
		-w /go/src/github.com/newrelic/nri-config-generator \
		-e SNYK_TOKEN \
		-e GO111MODULE=auto \
		snyk/snyk:golang snyk test --severity-threshold=high

.PHONY : ci/validate
ci/validate: ci/deps
	@docker run --rm -t \
			--name "nri-$(INTEGRATION)-validate" \
			-v $(CURDIR):/go/src/github.com/newrelic/nri-$(INTEGRATION) \
			-w /go/src/github.com/newrelic/nri-$(INTEGRATION) \
			$(BUILDER_TAG) make validate

.PHONY : ci/test
ci/test: ci/deps
	@docker run --rm -t \
		--name "nri-$(INTEGRATION)-test" \
		-v $(CURDIR):/go/src/github.com/newrelic/nri-$(INTEGRATION) \
		-w /go/src/github.com/newrelic/nri-$(INTEGRATION) \
		$(BUILDER_TAG) make test

.PHONY : ci/build
ci/build: ci/deps
ifdef TAG
	@docker run --rm -t \
		--name "nri-$(INTEGRATION)-build" \
		-v $(CURDIR):/go/src/github.com/newrelic/nri-$(INTEGRATION) \
		-w /go/src/github.com/newrelic/nri-$(INTEGRATION) \
		-e INTEGRATION \
		-e TAG
else
	@echo "===> $(INTEGRATION) ===  [ci/build] TAG env variable expected to be set"
	exit 1
endif
