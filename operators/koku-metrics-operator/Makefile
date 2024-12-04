PREVIOUS_VERSION ?= 3.3.1
VERSION ?= 3.3.2
PULLSPEC ?= quay.io/community-operator-pipeline-prod/koku-metrics-operator@sha256:9114f72f6adca60e18616786019d680883bb1c1dc88d317e7adeeec787054469

PDW=$(shell pwd)
OPERATOR_NAME=koku-metrics-operator
TOPDIR=$(abspath $(dir $(PWD))/../)
BINDIR=${TOPDIR}/bin

# Add the bin directory to the PATH
export PATH := $(BINDIR):$(PATH)
# A place to store the generated catalogs
CATALOG_DIR=${TOPDIR}/catalogs

# A place to store the operator catalog templates
OPERATOR_CATALOG_TEMPLATE_DIR = ${PDW}/catalog-templates

# A list of OCP versions to generate catalogs for
# This list can be customized to include the versions that are relevant to the operator
# DO NOT change this line (except for the versions) if you want to take advantage
# of the automated catalog promotion
OCP_VERSIONS=$(shell echo "v4.12 v4.13 v4.14 v4.15 v4.16 v4.17" )

.PHONY: catalogs
# replace this stub with one customized to serve your needs ... some examples below

# here are a few examples of different approaches to fulfilling this target
# comment out / customize the one that makes the most sense, or use them as examples in defining your own
#
# --- BASIC TEMPLATE ---
catalogs: basic validate-catalogs


# basic target provides an example FBC generation from a `basic` template type.
# this example takes a single file as input and generates a well-formed FBC operator contribution as an output
.PHONY: basic
basic: opm clean
	$(OPM) alpha render-template basic -o yaml ${OPERATOR_CATALOG_TEMPLATE_DIR}/template.yaml > ${OPERATOR_CATALOG_TEMPLATE_DIR}/catalog.yaml
	for version in $(OCP_VERSIONS); do \
		mkdir -p ${CATALOG_DIR}/$${version}/${OPERATOR_NAME}/ && \
		cp ${OPERATOR_CATALOG_TEMPLATE_DIR}/catalog.yaml ${CATALOG_DIR}/$${version}/${OPERATOR_NAME}/catalog.yaml; \
	done
	rm -f ${OPERATOR_CATALOG_TEMPLATE_DIR}/catalog.yaml

.PHONY: add-new-version
add-new-version:
	yq -i eval 'select(.schema == "olm.template.basic").entries[] |= select(.schema == "olm.channel").entries += [{"name" : "koku-metrics-operator.v$(VERSION)", "replaces": "koku-metrics-operator.v$(PREVIOUS_VERSION)"}]' ${OPERATOR_CATALOG_TEMPLATE_DIR}/template.yaml
	yq -i '.entries += [{"image": "$(PULLSPEC)", "schema": "olm.bundle"}]' ${OPERATOR_CATALOG_TEMPLATE_DIR}/template.yaml

# validate-catalogs target illustrates FBC validation
# all FBC must pass opm validation in order to be able to be used in a catalog
.PHONY: validate-catalogs
validate-catalogs: opm
	for version in $(OCP_VERSIONS); do \
		$(OPM) validate $(CATALOG_DIR)/$${version}/${OPERATOR_NAME} && echo "$${version} catalog validation passed" || echo "$${version} catalog validation failed"; \
	done

.PHONY: create-catalog-dir
create-catalog-dir:
	mkdir -p $(CATALOG_DIR)

.PHONY: clean
clean: create-catalog-dir
	find $(CATALOG_DIR) -type d -name ${OPERATOR_NAME} -exec rm -rf {} +

.PHONY: get-bundle
get-bundle:
	docker pull quay.io/community-operator-pipeline-prod/koku-metrics-operator:$(VERSION)
	docker inspect --format '{{.RepoDigests}}' quay.io/community-operator-pipeline-prod/koku-metrics-operator:$(VERSION)

OPM_VERSION=v1.48.0
OS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(shell uname -m | sed 's/x86_64/amd64/')

.PHONY: opm
OPM ?= $(BINDIR)/opm
opm: ## Download opm locally if necessary.
ifeq (,$(wildcard $(OPM)))
ifeq (, $(shell which opm 2>/dev/null))
	@{ \
	set -e ;\
	mkdir -p $(dir $(OPM)) ;\
	OS=$(shell go env GOOS) && ARCH=$(shell go env GOARCH) && \
	curl -Lo $(OPM) https://github.com/operator-framework/operator-registry/releases/download/$(OPM_VERSION)/$${OS}-$${ARCH}-opm ;\
	chmod +x $(OPM) ;\
	}
else
OPM = $(shell which opm)
endif
endif
