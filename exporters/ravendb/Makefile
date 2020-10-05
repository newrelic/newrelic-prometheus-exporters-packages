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

build: clone-repo
	@echo "=== $(NAME) === [ build ]:"
	$(WORK_DIR) go build -v -o $(BINS_DIR)/$(BINARY_NAME)

# Import fragments
include ../../scripts/package.mk

.PHONY: clone-repo build clean all