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

# A list of OCP versions to generate catalogs for
# This list can be customized to include the versions that are relevant to the operator
# DO NOT change this line (except for the versions) if you want to take advantage
# of the automated catalog promotion
OCP_VERSIONS=$(shell echo "v4.12 v4.13 v4.14 v4.15 v4.16" )


#
.PHONY: fbc-onboarding
fbc-onboarding: fbc-onboarding-deps ${BINDIR}/opm clean
	fbc-onboarding \
		--repo-root $(TOPDIR) \
		--operator-name $(OPERATOR_NAME) \
		--cache-dir $(TOPDIR)/.catalog_cache \
		--verbose

.PHONY: catalog
# replace this stub with one customized to serve your needs ... some examples below

# here are a few examples of different approaches to fulfilling this target
# comment out / customize the one that makes the most sense, or use them as examples in defining your own
#
# --- BASIC TEMPLATE ---
catalog: basic
#
# --- SEMVER TEMPLATE ---
#catalog: semver


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

#
# validate target illustrates FBC validation
# all FBC must pass opm validation in order to be able to be used in a catalog
.PHONY: validate
validate: ${BINDIR}/opm
	${BINDIR}/opm validate $(CATALOG_DIR) && echo "catalog validation passed" || echo "catalog validation failed"

.PHONY: create-catalog-dir
create-catalog-dir:
	mkdir -p $(CATALOG_DIR)

.PHONY: clean
clean: create-catalog-dir
	find $(CATALOG_DIR) -type d -name ${OPERATOR_NAME} -exec rm -rf {} +

# fbc-onboarding-deps installs the fbc-onboarding tool using pip
.PHONY: fbc-onboarding-deps
fbc-onboarding-deps:
	pip install git+https://github.com/redhat-openshift-ecosystem/operator-pipelines.git

OS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(shell uname -m | sed 's/x86_64/amd64/')

# Automatically download the opm binary
OPM_VERSION ?= v1.46.0
${BINDIR}/opm:
	if [ ! -d ${BINDIR} ]; then mkdir -p ${BINDIR}; fi
	curl -sLO https://github.com/operator-framework/operator-registry/releases/download/$(OPM_VERSION)/$(OS)-$(ARCH)-opm && chmod +x $(OS)-$(ARCH)-opm && mv $(OS)-$(ARCH)-opm ${BINDIR}/opm
