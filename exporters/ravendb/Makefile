# Exporter yml variables
NAME ?= ravendb
VERSION ?= 0.3.0
EXPORTER_REPO_URL ?= https://github.com/marcinbudny/ravendb_exporter
EXPORTER_HEAD ?= 0.3.0

IMPORT_PATH := $(subst https://,,$(EXPORTER_REPO_URL))
CLONED_REPO := $(GOPATH)/src/$(IMPORT_PATH)
WORK_DIR := cd $(CLONED_REPO) &&
BINARY_NAME := $(NAME)-exporter


TARGET_DIR = $(shell pwd)/target
BINS_DIR = $(TARGET_DIR)/bin/linux_amd64


all: clean clone-repo build package

clean: 
	@echo "=== $(NAME) === [ clean ]: Removing cloned folder"
	@if [ ! -d $(GOPATH) ]; then \
		echo "GOPATH is empty" ;\
		exit 1 ;\
	fi
	@CLONED_REPO="$(echo -e "${CLONED_REPO}" | tr -d '[:space:]')"
	@TARGET_DIR="$(echo -e "${TARGET_DIR}" | tr -d '[:space:]')"

	rm -rf $(CLONED_REPO) $(TARGET_DIR)

clone-repo:
	@echo "=== $(NAME) === [ clone-repo ]:"
	git clone $(EXPORTER_REPO_URL) $(CLONED_REPO)
	$(WORK_DIR) git checkout $(EXPORTER_HEAD)

build: clone-repo
	@echo "=== $(NAME) === [ build ]:"
	$(WORK_DIR) go build -v -o $(BINS_DIR)/$(BINARY_NAME)

# Import fragments
include ../../scripts/package.mk

.PHONY: clone-repo build clean all