# Exporter yml variables
NAME ?= githubactions
VERSION ?= 1.2.0
EXPORTER_REPO_URL ?= https://github.com/Spendesk/github-actions-exporter
EXPORTER_HEAD ?= v1.2

IMPORT_PATH := $(subst https://,,$(EXPORTER_REPO_URL))
CLONED_REPO := $(GOPATH)/src/$(IMPORT_PATH)
WORK_DIR := cd $(CLONED_REPO) &&
BINARY_NAME := $(NAME)-exporter

TARGET_DIR = $(shell pwd)/target
BINS_DIR = $(TARGET_DIR)/bin/linux_amd64


all: clean clone-repo get-definition-file build package

build: clone-repo deps-only
	@echo "=== $(NAME) === [ build ]:"
	$(WORK_DIR) go build -v -o $(BINS_DIR)/$(BINARY_NAME)

deps-only: clone-repo
	@echo "=== $(NAME) === [ deps ]: Installing package dependencies required by the project..."
	$(WORK_DIR) go mod download

# Import fragments
include ../../scripts/package.mk

.PHONY: all build clone-repo clean deps-only