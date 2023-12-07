# community-operator-pipeline

The community operator pipeline provides a platform for openshift community contributors to validate, share, and distribute openshift operators.

**community-hosted-pipeline**

The community hosted pipeline is used to certify the Operator bundles.
It’s an additional layer of validation that has to run within the Red Hat infrastructure. It is intended to be triggered upon the creation of a bundle pull request and successfully completes with merging it (configurable).


**community-release-pipeline**

The community release pipeline is responsible for releasing a bundle image which has passed certification.
It's intended to be triggered by the merge of a bundle pull request by the hosted pipeline.
It successfully completes once the bundle has been distributed to all relevant Operator catalogs
and appears in the Red Hat Ecosystem Catalog.

There are new github pull request labels introduced for the community operator pipeline.

| Label Name                               | Label Description                                                                           |
|------------------------------------------|---------------------------------------------------------------------------------------------|
| community-hosted-pipeline/started        | This label will be added to the PR when the hosted pipeline will be triggered.              |
| community-hosted-pipeline/passed         | This label will be added to the PR when the hosted pipeline will be passed successfully.    |
| community-hosted-pipeline/failed         | This label will be added to the PR when the hosted pipeline will be failed.                 |
| community-release-pipeline/started       | This label will be added to the PR when the release pipeline will be triggered.             |
| community-release-pipeline/passed        | This label will be added to the PR when the release pipeline will be passed successfully.   |
| community-release-pipeline/failed        | This label will be added to the PR when the release pipeline will be failed.                |
| validation-failed                        | This label will be added to the PR when the static tests will be failed in hosted pipeline. |



