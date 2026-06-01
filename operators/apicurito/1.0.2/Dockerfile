FROM scratch

# Core bundle labels.
LABEL operators.operatorframework.io.bundle.mediatype.v1=registry+v1
LABEL operators.operatorframework.io.bundle.manifests.v1=manifests/
LABEL operators.operatorframework.io.bundle.metadata.v1=metadata/
LABEL operators.operatorframework.io.bundle.package.v1=apicurito
LABEL operators.operatorframework.io.bundle.channels.v1=latest,candidate,stable,1.0.x
LABEL operators.operatorframework.io.bundle.channel.default.v1=latest
LABEL operators.operatorframework.io.metrics.builder=operator-sdk-v1.13.0+git
LABEL operators.operatorframework.io.metrics.mediatype.v1=metrics+v1
LABEL operators.operatorframework.io.metrics.project_layout=go.kubebuilder.io/v2

# Labels for testing.
LABEL operators.operatorframework.io.test.mediatype.v1=scorecard+v1
LABEL operators.operatorframework.io.test.config.v1=tests/scorecard/

# Copy files to locations specified by labels.
COPY manifests /manifests/
COPY metadata /metadata/
COPY tests/scorecard /tests/scorecard/

# Red Hat Openshift Labels
LABEL com.redhat.delivery.operator.bundle=true
LABEL com.redhat.openshift.versions="v4.6,v4.7,v4.8,v4.9"
