# Operator Publishing / Review settings

Each operator might have `ci.yaml` configuration file to be present in an operator directory (for example `community-operators/aqua/ci.yaml`). This configuration file is used by [community-operators](https://github.com/operator-framework/community-operators) pipeline to setup various features like `reviewers` or `operator versioning`.

> **Note:**
    One can create or modify `ci.yaml` file with a new operator version. This operation can be done in the same PR with other operator changes. 

## Reviewers

If you want to accelerate publishing your changes, consider adding yourself and others you trust to the `reviewers` list. If the author of PR will be in that list, changes she/he made will be taken as authorized changes. This will be the indicator for our pipeline that the PR is ready to merge automatically. 

> **Note:**
    If an author of PR is not in `reviewers` list or not in `ci.yaml` on `main` branch, PR will not be merged automatically. For security reasons, we are not checking reviewers from a fork.

> **Note:**
    If an auhor of PR is not in `reviewers` list and `reviewers` are present in `ci.yaml` file. All `reviewers` will be mentioned in PR comment to check for upcoming changes.

For this to work, it is required to setup reviewers in `ci.yaml` file. It can be done by adding `reviewers` tag with a list of GitHub usernames. For example

### Example
```
$ cat <path-to-operator>/ci.yaml
---
reviewers:
  - user1 
  - user2

```

## Operator versioning
Operators have multiple versions. When a new version is released, OLM can update an operator automatically. There are 2 update strategies possible, which are defined in `ci.yaml` at the operator top level.

### replaces-mode
Every next version defines which version will be replaced using `replaces` key in the CSV file. It means, that there is a possibility to omit some versions from the *update graph*. The best practice is to put them in a separate channel then.

### semver-mode
Every version will be replaced by the next higher version according to semantic versioning.

### Restrictions
A contributor can decide if `semver-mode` or `replaces-mode` mode will be used for a specific operator. By default, `replaces-mode` is activated, when `ci.yaml` file is present and contains `updateGraph: replaces-mode`. When a contributor decides to switch and use `semver-mode`, it will be specified in `ci.yaml` file or the key `updateGraph` will be missing.

### Example
```
$ cat <path-to-operator>/ci.yaml
---
# Use `replaces-mode` or `semver-mode`.
updateGraph: replaces-mode
```

## Kubernetes max version in CSV

Starting from kubernetes 1.22 some old APIs were deprecated ([Deprecated API Migration Guide from v1.22](https://kubernetes.io/docs/reference/using-api/deprecation-guide/#v1-22). Users can set `operatorhub.io/ui-metadata-max-k8s-version: "<version>"` in its CSV file to inform its maximum supported Kubernetes version. The following example will inform that operator can handle `1.21` as maximum Kubernetes version
```
$ cat <path-to-operators>/<name>/<version>/.../my.clusterserviceversion.yaml
metadata:
  annotations:
    operatorhub.io/ui-metadata-max-k8s-version: "1.21"
```
