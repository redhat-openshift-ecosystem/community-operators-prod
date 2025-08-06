# This Makefile provides a set of targets to generate and validate operator catalogs
# using the Operator Package Manager (opm) tool.

# The makefile should be placed in the root of the operator repository.
# for example at: <operator-repo>/operators/<operator-name>/Makefile

PWD=$(shell pwd)
OPERATOR_NAME=$(shell basename $(PWD))
TOPDIR=$(abspath $(dir $(PWD))/../)
BINDIR=${TOPDIR}/bin

# Add the bin directory to the PATH
export PATH := $(BINDIR):$(PATH)
# A place to store the generated catalogs
CATALOG_DIR=${TOPDIR}/catalogs

# The operator pipeline image to use for the fbc-onboarding target
OPERATOR_PIPELINE_IMAGE ?= quay.io/redhat-isv/operator-pipelines-images:released

# Define the paths for both auth files
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

# Catalogs are generated based on the mapping in ci.yaml file
# The example of the mapping is:
# fbc:
#   enabled: true
#   catalog_mapping:
#     - template_name: basic.yaml
#       catalog_names: ["v4.14", "v4.15", "v4.16"]
#       type: olm.template.basic
#     - template_name: semver.yaml
#       catalog_names: ["v4.13", "v4.17"]
#       type: olm.semver
catalogs: render_catalogs.sh ${BINDIR}/opm ${BINDIR}/yq clean
	@echo "Rendering catalogs from templates"
	@${BINDIR}/render_catalogs.sh


# validate-catalogs target illustrates FBC validation
# all FBC must pass opm validation in order to be able to be used in a catalog
.PHONY: validate-catalogs
validate-catalogs: ${BINDIR}/opm
	@find ${TOPDIR}/catalogs -type d -name ${OPERATOR_NAME} -exec \
		sh -c '${BINDIR}/opm validate $$(dirname "{}") && echo "✅ Catalog validation passed: {}" || echo "❌ Catalog validation failed: {}"' \;

.PHONY: create-catalog-dir
create-catalog-dir:
	@mkdir -p $(CATALOG_DIR)

.PHONY: clean
clean: create-catalog-dir
	@echo "Cleaning up the operator catalogs from $(CATALOG_DIR)"
	@find $(CATALOG_DIR) -type d -name ${OPERATOR_NAME} -exec rm -rf {} +


OS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(shell uname -m | sed 's/x86_64/amd64/')

# Automatically download the binaries and scripts and place them in the bin directory
OPM_VERSION ?= v1.46.0
${BINDIR}/opm:
	if [ ! -d ${BINDIR} ]; then mkdir -p ${BINDIR}; fi
	curl -sLO https://github.com/operator-framework/operator-registry/releases/download/$(OPM_VERSION)/$(OS)-$(ARCH)-opm && chmod +x $(OS)-$(ARCH)-opm && mv $(OS)-$(ARCH)-opm ${BINDIR}/opm

YQ_VERSION ?= v4.45.1
${BINDIR}/yq:
	if [ ! -d ${BINDIR} ]; then mkdir -p ${BINDIR}; fi
	curl -sLO https://github.com/mikefarah/yq/releases/download/$(YQ_VERSION)/yq_$(OS)_$(ARCH) && mv yq_$(OS)_$(ARCH) ${BINDIR}/yq && chmod +x ${BINDIR}/yq

.PHONY: render_catalogs.sh
render_catalogs.sh:
	@echo "Downloading render_catalogs.sh script"
	if [ ! -d ${BINDIR} ]; then mkdir -p ${BINDIR}; fi
	curl -sLO https://raw.githubusercontent.com/redhat-openshift-ecosystem/operator-pipelines/main/fbc/render_catalogs.sh
	mv render_catalogs.sh ${BINDIR}/render_catalogs.sh
	chmod +x ${BINDIR}/render_catalogs.sh

# push new version to community multi-nic-cni-operator
# update version-info.mk file with the new version before running this target
include version-info.mk

# A list of OCP versions to generate catalogs for
# This list can be customized to include the versions that are relevant to the operator
# DO NOT change this line (except for the versions) if you want to take advantage
# of the automated catalog promotion
OCP_VERSIONS=$(shell echo "v4.12 v4.13 v4.14 v4.15 v4.16 v4.17 v4.18 v4.19" )

.PHONY: push-bundle
push-bundle:
	if [ -z "${VERSION}" ]; then \
		echo "Error: Version is not set"; \
		exit 1; \
	fi; \
	if [ -z "${OLD_VERSION}" ]; then \
		echo "Error: OLD_VERSION is not set"; \
		exit 1; \
	fi; \
	if [ -z "${CHANNELS}" ]; then \
		echo "Error: CHANNELS is not set"; \
		exit 1; \
	fi; \
	# Checkout the git repository with release branch \
	git clone -b release-v${VERSION} https://github.com/foundation-model-stack/multi-nic-cni.git multi-nic-cni-operator && \
	# Change to the cloned directory \
	cd multi-nic-cni-operator && \
	# Build the opm binary \
	make opm && \
	# Tidy and vendor Go modules \
	go mod tidy && \
	go mod vendor && \
	# Execute the make bundle command \
	make bundle && \
	# Return to the original directory \
	cd ..; \
	# Check if the ${VERSION} directory exists and delete it \
	if [ -d "${VERSION}" ]; then \
		echo "Directory ${VERSION} exists and will be deleted"; \
		rm -rf ${VERSION}; \
	fi; \
	# Create the target directory for the specified version \
	mkdir -p ${VERSION}; \
	# Copy the subdirectories under the generated bundle directory to the target directory \
	cp -r multi-nic-cni-operator/bundle/* ${VERSION}; \
	# Copy the bundle.Dockerfile to the target directory \
	cp multi-nic-cni-operator/bundle.Dockerfile ${VERSION}; \
	# Create the release-config.yaml file \
	echo "---" > ${VERSION}/release-config.yaml; \
	echo "catalog_templates:" >> ${VERSION}/release-config.yaml; \
	for version in ${OCP_VERSIONS}; do \
		echo "  - template_name: $${version}.yaml" >> ${VERSION}/release-config.yaml; \
		echo "    channels: [${CHANNELS}]" >> ${VERSION}/release-config.yaml; \
		echo "    replaces: multi-nic-cni-operator.v${OLD_VERSION}" >> ${VERSION}/release-config.yaml; \
	done; \
	# Remove the cloned git repository \
	rm -rf multi-nic-cni-operator
