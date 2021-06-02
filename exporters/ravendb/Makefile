# Exporter yml variables
NAME ?= ravendb
VERSION ?= 0.3.0
EXPORTER_REPO_URL ?= https://github.com/marcinbudny/ravendb_exporter
EXPORTER_HEAD ?= 0.3.0
EXPORTER_DEFAULT_PORT ?= 9440
IMPORT_PATH := $(subst https://,,$(EXPORTER_REPO_URL))
CLONED_REPO := $(GOPATH)/src/$(IMPORT_PATH)
WORK_DIR := cd $(CLONED_REPO) &&
BINARY_NAME := $(NAME)-exporter
NRI_CONFIG_GENERATOR_DIR := cd ../../nri-config-generator &&
GOFLAGS          = -mod=readonly

TARGET_DIR = $(shell pwd)/target
BINS_DIR = $(TARGET_DIR)/bin/linux_amd64


all: clean clone-repo get-definition-file build package

build: clone-repo
	@echo "=== $(NAME) === [ build ]:"
	$(WORK_DIR) GO111MODULE=off go build -v -o $(BINS_DIR)/$(BINARY_NAME)
	cp $(NAME).json.tmpl ../../nri-config-generator/templates/$(NAME).json.tmpl
	$(NRI_CONFIG_GENERATOR_DIR) go build -ldflags "-X main.integration=$(INTEGRATION_NAME) -X main.defPort=$(EXPORTER_DEFAULT_PORT)" -o $(BINS_DIR)/nri-$(NAME) github.com/newrelic/nri-config-generator

# Import fragments
include ../../scripts/package.mk

.PHONY: clone-repo build clean all