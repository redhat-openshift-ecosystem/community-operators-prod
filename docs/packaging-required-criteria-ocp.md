## OKD/OpenShift Catalogs criteria and options

### Overview

To distribute on OpenShift Catalogs, you will need to comply with the same standard criteria defined for `OperatorHub.io` (see [Common recommendations and suggestions](https://olm.operatorframework.io/docs/best-practices/common/#validate-your-bundle-before-publish-it)). Then, additionally, you have some requirements and options which follow.

> **IMPORTANT** Kubernetes has been deprecating API(s) which will be removed and no longer available in `1.22` and in the Openshift version `4.9`. Note that your project will be unable to use them on `OCP 4.9/K8s 1.22` and then, it is strongly recommended to check [Deprecated API Migration Guide from v1.22][k8s-deprecated-guide] and ensure that your projects have them migrated and are not using any deprecated API.

> Note that your operator using them will not work in  `1.22` and in the Openshift version `4.9`. OpenShift 4.8 introduces two new alerts that fire when an API that will be removed in the next release is in use. Check the event alerts of your Operators running on 4.8 and ensure that you will not find any warning about these API(s) still being used by it.

> Also, to prevent workflow issues, its users will need to have installed in their OCP cluster a version of your operator compatible with 4.9 before they try to upgrade their cluster from any previous version to 4.9 or higher. In this way, it is recommended to ensure that your operators are no longer using these API(s) versions. However, If you still need to publish the operator bundles with any of these API(s) for use on earlier k8s/OCP versions, ensure that the operator bundle is configured accordingly.

> Taking the actions below will help prevent users from installing versions of your operator on an incompatible version of OCP, and also prevent them from upgrading to a newer version of OCP that would be incompatible with the version of your operator that is currently installed on their cluster.

### Configure the max OpenShift Version compatible

Use the `olm.openShiftMaxVersion` annotation in the CSV to prevent the user from upgrading their OCP cluster before upgrading the installed operator version to any distribution which is compatible with:

```yaml
apiVersion: operators.coreos.com/v1alpha1
kind: ClusterServiceVersion
metadata:
  annotations:
    # Prevent cluster upgrades to OpenShift Version 4.9 when this
    # bundle is installed on the cluster
    "olm.properties": '[{"type": "olm.maxOpenShiftVersion", "value": "4.8"}]'
```

The CSV annotation will eventually prevent the user from upgrading their OCP cluster before they have installed a version of your operator which is compatible with `4.9`. However, note that it is important to make these changes now as users running workloads with deprecated API(s) that are looking to upgrade to OCP 4.9 will need to be running operators that have this annotation set in order to prevent the cluster upgrade and potentially adversely impacting their crucial workloads.

This option is useful when you know that the current version of your project will not work well on some specific Openshift version.

### Configure the Openshift distribution

Use the annotation `com.redhat.openshift.versions` in `bundle/metadata/annotations.yaml` to ensure that the index image will be generated with its OCP Label, to prevent the bundle from being distributed on to 4.9:

```
com.redhat.openshift.versions: "v4.6-v4.8"
```

This option is also useful when you know that the current version of your project will not work well on some specific OpenShift version. By using it you defined the Openshift versions where the Operator should be distributed and the Operator will not appear in a catalog of an Openshift version that is outside of the range. You must use it if you are distributing a solution that contains deprecated API(s) and will no longer be available in later versions. For more information see [Managing OpenShift Versions][managing-openshift-versions].

### Validate the bundle with the common criteria to distribute via OLM with SDK

Also, you can check the bundle via [`operator-sdk bundle validate`][sdk-cli-bundle-validate] against the suite  Validator [Community Operators][optional-validators] and the K8s Version that you are intended to publish:

```sh
operator-sdk bundle validate ./bundle --select-optional suite=operatorframework --optional-values=k8s-version=1.22
```
**NOTE:** The validators only check the manifests which are shipped in the bundle. They are unable to ensure that the project's code does not use the [Deprecated/Removed API(s) in 1.22][k8s-deprecated-guide] and/or that it does not have as dependency another operator that uses them.

### Validate the bundle with the specific criteria to distribute in Openshift catalogs

**Pre-requirement**
Download the [binary](https://github.com/redhat-openshift-ecosystem/ocp-olm-catalog-validator/releases). You might want to keep it in your `$GOPTH/bin`

Then, we can use the experimental [OpenShift OLM Catalog Validator](https://github.com/redhat-openshift-ecosystem/ocp-olm-catalog-validator) to check your Operator bundle.
In this case, we need to inform the bundle and the annotations.yaml file paths:

```
$ ocp-olm-catalog-validator my-bundle-path/bundle  --optional-values="file=bundle-path/bundle/metadata/annotations.yaml"
```

Following is an example of an Operator bundle that uses the removed APIs in 1.22 and is not configured accordingly:

```sh
$ ocp-olm-catalog-validator bundle/ --optional-values="file=bundle/metadata/annotations.yaml"
WARN[0000] Warning: Value memcached-operator.v0.0.1: this bundle is using APIs which were deprecated and removed in v1.22. More info: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22. Migrate the API(s) for CRD: (["memcacheds.cache.example.com"])
ERRO[0000] Error: Value : (memcached-operator.v0.0.1) olm.maxOpenShiftVersion csv.Annotations not specified with an OCP version lower than 4.9. This annotation is required to prevent the user from upgrading their OCP cluster before they have installed a version of their operator which is compatible with 4.9. For further information see https://docs.openshift.com/container-platform/4.8/operators/operator_sdk/osdk-working-bundle-images.html#osdk-control-compat_osdk-working-bundle-images
ERRO[0000] Error: Value : (memcached-operator.v0.0.1) this bundle is using APIs which were deprecated and removed in v1.22. More info: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22. Migrate the APIs for this bundle is using APIs which were deprecated and removed in v1.22. More info: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22. Migrate the API(s) for CRD: (["memcacheds.cache.example.com"]) or provide compatible version(s) via the labels. (e.g. LABEL com.redhat.openshift.versions='4.6-4.8')
```

[sdk-cli-bundle-validate]: https://sdk.operatorframework.io/docs/cli/operator-sdk_bundle_validate/
[managing-openshift-versions]: https://redhat-connect.gitbook.io/certified-operator-guide/ocp-deployment/operator-metadata/bundle-directory/managing-openshift-versions
[optional-validators]: https://olm.operatorframework.io/docs/tasks/creating-operator-bundle/#optional-validation
[k8s-deprecated-guide]: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22