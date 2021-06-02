INTEGRATION_NAME=
EXPORTER_DEFAULT_PORT=
compile:
ifeq ($(INTEGRATION_NAME),)
$(error Missing required argument INTEGRATION_NAME.)
endif
ifeq ($(EXPORTER_DEFAULT_PORT),)
$(error Missing required argument EXPORTER_DEFAULT_PORT.)
endif
	go build -ldflags "-X main.integration=$(INTEGRATION_NAME) -X main.defPort=$(EXPORTER_DEFAULT_PORT)" -o dist/nri-$(INTEGRATION_NAME) github.com/newrelic/nri-config-generator
	echo "created executable nri-$(INTEGRATION_NAME)"