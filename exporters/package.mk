PROJECT_NAME       = $(NAME)
PACKAGE_TYPES     ?= deb rpm tarball
SOURCE_DIR         = $(TARGET_DIR)/source
PACKAGES_DIR       = $(TARGET_DIR)/packages
TARBALL_DIR       ?= $(PACKAGES_DIR)/tarball
PKG_TARBALL       ?= true
GOARCH            ?= amd64
VERSION           ?= 0.0.0
RELEASE           ?= $(VERSION)

LICENSE            = "https://newrelic.com/terms (also see LICENSE.txt installed with this package)"
VENDOR             = "New Relic, Inc."
PACKAGER           = "New Relic, Inc."
PACKAGE_URL        = "https://github.com/newrelic/newrelic-prometheus-exporters-packages"
SUMMARY            = "Prometheus exporter for $(NAME) $(EXPORTER_REPO_URL) "
DESCRIPTION        = "Prometheus exporters help exporting existing metrics from third-party systems as Prometheus metrics."

FPM_COMMON_OPTIONS = --verbose -C $(SOURCE_DIR) -s dir -n $(PROJECT_NAME) -v $(VERSION) --iteration $(RELEASE) --prefix "" --license $(LICENSE) --vendor $(VENDOR) -m $(PACKAGER) --url $(PACKAGE_URL) --config-files /etc/newrelic-infra/ --description "$$(printf $(DESCRIPTION))"
FPM_DEB_OPTIONS    = -t deb -p $(PACKAGES_DIR)/deb/
FPM_RPM_OPTIONS    = -t rpm -p $(PACKAGES_DIR)/rpm/ --epoch 0 --rpm-summary $(SUMMARY)

package: $(PACKAGE_TYPES)


prep-pkg-env:
	@if [ ! -d $(BINS_DIR) ]; then \
		echo "=== Main === [ prep-pkg-env ]: no built binaries found. Run 'make create-bins'" ;\
		exit 1 ;\
	fi
	@echo "=== Main === [ prep-pkg-env ]: preparing a clean packaging environment..."
	@rm -rf $(SOURCE_DIR)
	@mkdir -p $(SOURCE_DIR)/usr/local/prometheus-exporters/$(NAME) $(SOURCE_DIR)/etc/newrelic-infra/integrations.d
	@echo "=== Main === [ prep-pkg-env ]: adding built binaries and configuration and definition files..."
	@cp $(BINS_DIR)/* $(SOURCE_DIR)/usr/local/prometheus-exporters/$(NAME)
	@chmod 755 $(SOURCE_DIR)/usr/local/prometheus-exporters/$(NAME)/*
	@echo "=== Main === [ prep-pkg-env ]: adding license..."
	@cp LICENSE $(SOURCE_DIR)/usr/local/prometheus-exporters/$(NAME)

deb: prep-pkg-env
	@echo "=== Main === [ deb ]: building DEB package..."
	@mkdir -p $(PACKAGES_DIR)/deb
	@fpm $(FPM_COMMON_OPTIONS) $(FPM_DEB_OPTIONS) .

rpm: prep-pkg-env
	@echo "=== Main === [ rpm ]: building RPM package..."
	@mkdir -p $(PACKAGES_DIR)/rpm
	@fpm $(FPM_COMMON_OPTIONS) $(FPM_RPM_OPTIONS) .

FILENAME_TARBALL_LINUX = $(PROJECT_NAME)_linux_$(VERSION)_$(GOARCH).tar.gz
tarball: prep-pkg-env
	@echo "=== Main === [ tar ]: building Tarball package..."
	@mkdir -p $(TARBALL_DIR)
	tar -czf $(TARBALL_DIR)/$(FILENAME_TARBALL_LINUX) -C $(SOURCE_DIR) ./

.PHONY: package prep-pkg-env $(PACKAGE_TYPES)
