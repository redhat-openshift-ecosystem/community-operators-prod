Thanks for submitting your Operator. Please check the below list before you create your Pull Request.

### New Submissions

* [ ] Are you familiar with our [contribution guidelines](https://redhat-openshift-ecosystem.github.io/operator-pipelines/users/contributing-prerequisites/)?
* [ ] Are you familiar with our [operator pipeline](https://redhat-openshift-ecosystem.github.io/operator-pipelines/users/pipelines_overview/)?
* [ ] Have you tested your Operator with all Custom Resource Definitions [packaging](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/)?
* [ ] Have you tested your Operator in all supported [installation modes](https://github.com/operator-framework/operator-lifecycle-manager/blob/master/doc/design/building-your-csv.md#operator-metadata)?
* [ ] Have you considered whether you want to use [semantic versioning order](https://github.com/operator-framework/community-operators/blob/master/docs/operator-ci-yaml.md#semver-mode)?
* [ ] Is your submission [signed](https://redhat-openshift-ecosystem.github.io/operator-pipelines/users/contributing-prerequisites/#sign-your-work)?
* [ ] Is operator [icon](https://github.com/operator-framework/community-operators/blob/master/docs/packaging-operator.md#operator-icon) set?

### Your submission should not

* [ ] Add more than one operator bundle per PR
* [ ] Modify any operator
* [ ] Rename an operator
* [ ] Modify any files outside the above-mentioned folders
* [ ] Contain more than one commit. **Please squash your commits.**

### Operator Description must contain (in order)

1. [ ] Description of the managed Application and where to find more information
2. [ ] Features and capabilities of your Operator and how to use it
3. [ ] Any manual steps about potential pre-requisites for using your Operator

### Operator Metadata should contain

* [ ] Human-readable name and 1-liner description about your Operator
* [ ] Valid [category](https://github.com/operator-framework/community-operators/blob/master/docs/packaging-operator.md#categories) name <sup>1</sup>
* [ ] One of the pre-defined [capability levels](https://github.com/operator-framework/operator-courier/blob/4d1a25d2c8d52f7de6297ec18d8afd6521236aa2/operatorcourier/validate.py#L556)<sup>2</sup>
* [ ] Links to the maintainer, source code and documentation
* [ ] Example templates for all Custom Resource Definitions intended to be used
* [ ] A quadratic logo

Remember that you can preview your CSV [here](https://operatorhub.io/preview).

--

<sup>1</sup> If you feel your Operator does not fit any of the pre-defined categories, file an issue against this repo and explain your need

<sup>2</sup> For more information see [here](https://sdk.operatorframework.io/docs/overview/#operator-capability-level)
