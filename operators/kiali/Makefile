# This Makefile provides a set of targets to generate and validate operator catalogs
# using the Operator Package Manager (opm) tool.

# The makefile should be placed in the root of the operator repository.
# for example at: <operator-repo>/operators/<operator-name>/Makefile

# A user can customize "catalog" target to generate the operator catalog in a way
# that suits the operator.
# OPM allows for the generation of catalogs using different templates.
# - basic: generates a basic catalog
# - semver: generates a catalog with semver versioning

PDW=$(shell pwd)
OPERATOR_NAME=$(shell basename $(PDW))
TOPDIR=$(abspath $(dir $(PWD))/../)
BINDIR=${TOPDIR}/bin

# Add the bin directory to the PATH
export PATH := $(BINDIR):$(PATH)
# A place to store the generated catalogs
CATALOG_DIR=${TOPDIR}/catalogs

# A place to store the operator catalog templates
OPERATOR_CATALOG_TEMPLATE_DIR = ${PDW}/catalog-templates

# The operator pipeline image to use for the fbc-onboarding target
OPERATOR_PIPELINE_IMAGE ?= quay.io/redhat-isv/operator-pipelines-images:released

# Define the paths for both auth files
DOCKER_CONFIG := $(HOME)/.docker/config.json
CONTAINERS_AUTH := $(XDG_RUNTIME_DIR)/containers/auth.json

# A list of OCP versions to generate catalogs for
# This list can be customized to include the versions that are relevant to the operator
# DO NOT change this line (except for the versions) if you want to take advantage
# of the automated catalog promotion
OCP_VERSIONS=$(shell echo "v4.12 v4.13 v4.14 v4.15 v4.16 v4.17" )


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
	podman run \
		--rm \
		--user $(id -u):$(id -g) \
		--security-opt label=disable \
		--pull always \
		-v $(TOPDIR):/workspace \
		$$CONFIG_VOLUME \
		$(OPERATOR_PIPELINE_IMAGE) fbc-onboarding \
			--repo-root /workspace \
			--operator-name $(OPERATOR_NAME) \
			--cache-dir /workspace/.catalog_cache

.PHONY: catalogs
# replace this stub with one customized to serve your needs ... some examples below

# here are a few examples of different approaches to fulfilling this target
# comment out / customize the one that makes the most sense, or use them as examples in defining your own
#
# --- BASIC TEMPLATE ---
catalogs: basic
#
# --- SEMVER TEMPLATE ---
#catalogs: semver


# basic target provides an example FBC generation from a `basic` template type.
# this example takes a single file as input and generates a well-formed FBC operator contribution as an output
.PHONY: basic
basic: ${BINDIR}/opm clean
	for version in $(OCP_VERSIONS); do \
		mkdir -p ${CATALOG_DIR}/$${version}/${OPERATOR_NAME}/ && \
		${BINDIR}/opm alpha render-template basic -o yaml ${OPERATOR_CATALOG_TEMPLATE_DIR}/$${version}.yaml > ${CATALOG_DIR}/$${version}/${OPERATOR_NAME}/catalog.yaml; \
	done


# semver target provides an example FBC generation from a `semver` template type.
# this example takes a single file as input and generates a well-formed FBC operator contribution as an output
.PHONY: semver
semver: ${BINDIR}/opm clean
	for version in $(OCP_VERSIONS); do \
		mkdir -p ${CATALOG_DIR}/$${version}/${OPERATOR_NAME}/ && \
		${BINDIR}/opm alpha render-template semver -o yaml  ${OPERATOR_CATALOG_TEMPLATE_DIR}/$${version}.yaml > ${CATALOG_DIR}/$${version}/${OPERATOR_NAME}/catalog.yaml; \
	done


# validate-catalogs target illustrates FBC validation
# all FBC must pass opm validation in order to be able to be used in a catalog
.PHONY: validate-catalogs
validate-catalogs: ${BINDIR}/opm
	for version in $(OCP_VERSIONS); do \
		${BINDIR}/opm validate $(CATALOG_DIR)/$${version}/${OPERATOR_NAME} && echo "$${version} catalog validation passed" || echo "$${version} catalog validation failed"; \
	done

.PHONY: create-catalog-dir
create-catalog-dir:
	mkdir -p $(CATALOG_DIR)

.PHONY: clean
clean: create-catalog-dir
	find $(CATALOG_DIR) -type d -name ${OPERATOR_NAME} -exec rm -rf {} +


OS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(shell uname -m | sed 's/x86_64/amd64/')

# Automatically download the opm binary
OPM_VERSION ?= v1.46.0
${BINDIR}/opm:
	if [ ! -d ${BINDIR} ]; then mkdir -p ${BINDIR}; fi
	curl -sLO https://github.com/operator-framework/operator-registry/releases/download/$(OPM_VERSION)/$(OS)-$(ARCH)-opm && chmod +x $(OS)-$(ARCH)-opm && mv $(OS)-$(ARCH)-opm ${BINDIR}/opm
