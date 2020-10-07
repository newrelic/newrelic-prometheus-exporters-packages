# Exporter yml variables
NAME ?= github
VERSION ?= 2.0.1
EXPORTER_REPO_URL ?= https://github.com/infinityworks/github-exporter
EXPORTER_HEAD ?= 2.0.1-alpha

IMPORT_PATH := $(subst https://,,$(EXPORTER_REPO_URL))
CLONED_REPO := $(GOPATH)/src/$(IMPORT_PATH)
WORK_DIR := cd $(CLONED_REPO) &&
BINARY_NAME := $(NAME)-exporter

TARGET_DIR = $(shell pwd)/target
BINS_DIR = $(TARGET_DIR)/bin/linux_amd64


all: clean clone-repo build package

build: clone-repo deps-only
	@echo "=== $(NAME) === [ build ]:"
	$(WORK_DIR) go build -v -o $(BINS_DIR)/$(BINARY_NAME)

deps-only: clone-repo
	@echo "=== $(NAME) === [ deps ]: Installing package dependencies required by the project..."
	$(WORK_DIR) go mod download

# Import fragments
include ../../scripts/package.mk

.PHONY: all build clone-repo clean deps-only