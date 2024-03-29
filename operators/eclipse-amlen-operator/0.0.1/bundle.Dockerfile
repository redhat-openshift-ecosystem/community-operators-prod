FROM scratch

# Core bundle labels.
LABEL operators.operatorframework.io.bundle.mediatype.v1=registry+v1
LABEL operators.operatorframework.io.bundle.manifests.v1=manifests/
LABEL operators.operatorframework.io.bundle.metadata.v1=metadata/
LABEL operators.operatorframework.io.bundle.package.v1=eclipse-amlen-operator
LABEL operators.operatorframework.io.bundle.channels.v1=preview
LABEL operators.operatorframework.io.bundle.channel.default.v1=stable

# Copy files to locations specified by labels.
COPY eclipse-amlen-operator/0.0.1/manifests /manifests/
COPY eclipse-amlen-operator/0.0.1/metadata /metadata/
