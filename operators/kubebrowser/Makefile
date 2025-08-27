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

BROWSER_CRD_DESCRIPTION=Deploys a browser instance
BROWSERSYSTEM_CRD_DESCRIPTION=Configures resources needed for running Browsers


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
catalogs: render_catalogs.sh ${BINDIR}/opm ${BINDIR}/yq clean
	@echo "Rendering catalogs from templates"
	@${BINDIR}/render_catalogs.sh
# opm render doesn't include description on olm.gvk 
	CATALOGS=$$(yq -r '.fbc.catalog_mapping[].catalog_names[]' ci.yaml | awk '!seen[$$0]++'); \
	for name in $$CATALOGS; do \
		echo "Inserting defaults for $$name"; \
		yq '.properties[] |= (select(.type=="olm.gvk" and .value.kind=="Browser").value.description="${BROWSER_CRD_DESCRIPTION}") |= (select(.type=="olm.gvk" and .value.kind=="BrowserSystem").value.description="${BROWSERSYSTEM_CRD_DESCRIPTION}")' -i ../../catalogs/$$name/kubebrowser/catalog.yaml; \
	done


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
	if [ ! -d ${BINDIR} ]; then mkdir -p ${BINDIR}; fi
	curl -sLO https://raw.githubusercontent.com/redhat-openshift-ecosystem/operator-pipelines/main/fbc/render_catalogs.sh
	mv render_catalogs.sh ${BINDIR}/render_catalogs.sh
	chmod +x ${BINDIR}/render_catalogs.sh
