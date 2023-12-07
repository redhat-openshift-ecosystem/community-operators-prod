# Openshift Community Operators
[![License](http://img.shields.io/:license-apache-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)

## About this repository

This repo is the canonical source for Kubernetes Operators that appear on [OpenShift Container Platform](https://openshift.com) and [OKD](https://www.okd.io/).

**NOTE** The index catalog `registry.redhat.io/redhat/redhat-operator-index:v<OCP Version>` is built from this repository and it is
consumed by Openshift and OKD to create their sources and built their catalog. To know more about how
Openshift catalog are built see the [documentation](https://docs.openshift.com/container-platform/4.6/operators/understanding/olm-rh-catalogs.html#olm-rh-catalogs_olm-rh-catalogs).

## Testing and release process for the Operator

Refer to the [Community operator pipeline documentation](docs/community-operator-pipeline.md) .

## Documentation

Full documentation is generated via [mkdoc](https://www.mkdocs.org/) and is located at [https://redhat-openshift-ecosystem.github.io/community-operators-prod/](https://redhat-openshift-ecosystem.github.io/community-operators-prod/)

## IMPORTANT NOTICE

Some APIs versions are deprecated and are OR will no longer be served on the Kubernetes version
`1.22/1.25/1.26` and consequently on vendors like Openshift `4.9/4.12/4.13`.

**What does it mean for you?**

Operator bundle versions using the removed APIs can not work successfully from the respective releases.
Therefore, it is recommended to check if your solutions are failing in these scenarios to stop using these versions
OR by setting the `"olm.properties": '[{"type": "olm.maxOpenShiftVersion", "value": "<OCP version>"}]'`
to block cluster admins upgrades when they have Operator versions installed that can **not**
work well in OCP versions higher than the value informed. Also, by defining a valid OCP range via the annotation `com.redhat.openshift.versions`
into the `metadata/annotations.yaml` for our solution does **not** end up shipped on OCP/OKD versions where it cannot be installed.

> WARNING: `olm.maxOpenShiftVersion` should ONLY be used if you are 100% sure that your Operator bundle version
> cannot work in upper releases. Otherwise, you might provide a bad user experience. Be aware that cluster admins
> will be unable to upgrade their clusters with your solution installed. Then, suppose you do not provide any upper
> version and a valid upgrade path for those who have your Operator installed be able to upgrade it and consequently
> be allowed to upgrade their cluster version (i.e from OCP 4.10 to 4.11). In that case, cluster admins might
> choose to uninstall your Operator and no longer use it so that they can move forward and upgrade their cluster
> version without it.

Please, make sure you check the following announcements:
- [How to deal with removal of v1beta1 CRD removals in Kubernetes 1.22 / OpenShift 4.9](https://github.com/redhat-openshift-ecosystem/community-operators-prod/discussions/138)
- [Kubernetes API removals on 1.25/1.26 and Openshift 4.12/4.13 might impact your Operator. How to deal with it?](https://github.com/redhat-openshift-ecosystem/community-operators-prod/discussions/1182)

## Reporting Bugs

Use the issue tracker in this repository to report bugs.

[k8s-deprecated-guide]: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22
