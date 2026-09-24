# This Makefile provides targets to migrate and maintain the operator catalog
# using the File-Based Catalog workflow.

PWD := $(shell pwd)
OPERATOR_NAME := $(shell basename $(PWD))
TOPDIR := $(abspath $(dir $(PWD))/../)
BINDIR := $(TOPDIR)/bin
CATALOG_DIR := $(TOPDIR)/catalogs

export PATH := $(BINDIR):$(PATH)

OPERATOR_PIPELINE_IMAGE ?= quay.io/redhat-isv/operator-pipelines-images:released
DOCKER_CONFIG := $(HOME)/.docker/config.json
CONTAINERS_AUTH := $(XDG_RUNTIME_DIR)/containers/auth.json

.PHONY: fbc-onboarding
fbc-onboarding: clean
	@if [ -f $(DOCKER_CONFIG) ]; then \
		echo "Using Docker config file: $(DOCKER_CONFIG)"; \
		CONFIG_VOLUME="-v $(DOCKER_CONFIG):/root/.docker/config.json"; \
	elif [ -f $(CONTAINERS_AUTH) ]; then \
		echo "Using containers auth file: $(CONTAINERS_AUTH)"; \
		CONFIG_VOLUME="-v $(CONTAINERS_AUTH):/root/.docker/config.json"; \
	else \
		echo "No authentication file found."; \
	fi; \
	podman run --rm \
		--user $$(id -u):$$(id -g) \
		--security-opt label=disable \
		--pull always \
		-v $(TOPDIR):/workspace \
		$$CONFIG_VOLUME \
		$(OPERATOR_PIPELINE_IMAGE) fbc-onboarding \
			--repo-root /workspace \
			--operator-name $(OPERATOR_NAME) \
			--cache-dir /workspace/.catalog_cache

.PHONY: catalogs
catalogs: render_catalogs.sh $(BINDIR)/opm $(BINDIR)/yq clean
	@echo "Rendering catalogs from templates"
	@$(BINDIR)/render_catalogs.sh

.PHONY: validate-catalogs
validate-catalogs: $(BINDIR)/opm
	@find $(CATALOG_DIR) -type d -name $(OPERATOR_NAME) -exec \
		sh -c '$(BINDIR)/opm validate $$(dirname "{}") && echo "Catalog validation passed: {}"' \;

.PHONY: create-catalog-dir
create-catalog-dir:
	@mkdir -p $(CATALOG_DIR)

.PHONY: clean
clean: create-catalog-dir
	@find $(CATALOG_DIR) -type d -name $(OPERATOR_NAME) -exec rm -rf {} +

OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m | sed 's/x86_64/amd64/')

OPM_VERSION ?= v1.46.0
$(BINDIR)/opm:
	@mkdir -p $(BINDIR)
	curl -fsSLO https://github.com/operator-framework/operator-registry/releases/download/$(OPM_VERSION)/$(OS)-$(ARCH)-opm
	chmod +x $(OS)-$(ARCH)-opm
	mv $(OS)-$(ARCH)-opm $(BINDIR)/opm

YQ_VERSION ?= v4.45.1
$(BINDIR)/yq:
	@mkdir -p $(BINDIR)
	curl -fsSLO https://github.com/mikefarah/yq/releases/download/$(YQ_VERSION)/yq_$(OS)_$(ARCH)
	mv yq_$(OS)_$(ARCH) $(BINDIR)/yq
	chmod +x $(BINDIR)/yq

.PHONY: render_catalogs.sh
render_catalogs.sh:
	@mkdir -p $(BINDIR)
	curl -fsSLO https://raw.githubusercontent.com/redhat-openshift-ecosystem/operator-pipelines/main/fbc/render_catalogs.sh
	mv render_catalogs.sh $(BINDIR)/render_catalogs.sh
	chmod +x $(BINDIR)/render_catalogs.sh
