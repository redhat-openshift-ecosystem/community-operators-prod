FROM scratch

LABEL operators.operatorframework.io.bundle.mediatype.v1=registry+v1
LABEL operators.operatorframework.io.bundle.manifests.v1=manifests/
LABEL operators.operatorframework.io.bundle.metadata.v1=metadata/
LABEL operators.operatorframework.io.bundle.package.v1=3scale-community-operator
LABEL operators.operatorframework.io.bundle.channels.v1=threescale-2.10,threescale-2.6,threescale-2.7,threescale-2.8,threescale-2.9
LABEL operators.operatorframework.io.bundle.channel.default.v1=threescale-2.10
LABEL com.redhat.openshift.versions=v4.6-v4.8

COPY manifests /manifests/
COPY metadata /metadata/
