# Openshift Community Operators
[![License](http://img.shields.io/:license-apache-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)

## About this repository

This repo is the canonical source for Kubernetes Operators that appear on [OpenShift Container Platform](https://openshift.com)

## Documentation

Full documentation is generated via [mkdoc](https://www.mkdocs.org/) and is located at [https://redhat-openshift-ecosystem.github.io/community-operators-prod/](https://redhat-openshift-ecosystem.github.io/community-operators-prod/)

## IMPORTANT NOTICE

**IMPORTANT** Kubernetes has been deprecating API(s) which will be removed and no longer available in `1.22` and in the Openshift version `4.9`. Note that your project will be unable to use them on `OCP 4.9/K8s 1.22` and then, it is strongly recommended to check [Deprecated API Migration Guide from v1.22][k8s-deprecated-guide] and ensure that your projects have them migrated and are not using any deprecated API.

### FOR OPENSHIFT COMMUNITY OPERATORS

However, If you still need to publish the operator bundles with any of these API(s) on [OpenShift Container Platform](https://openshift.com) and [OKD](https://okd.io) ensure that it is configured with the criteria defined in [OKD/OpenShift Catalogs criteria and options](https://operator-framework.github.io/community-operators/packaging-required-criteria-ocp/).

## Reporting Bugs

Use the issue tracker in this repository to report bugs.

[k8s-deprecated-guide]: https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22
