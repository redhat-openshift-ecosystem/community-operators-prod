# Updating existing Operators

## Updating operator version

> **Note:**
    It is strongly recommended to bump the operator version if possible.

Operator version update can be done by creating a new directory with `version` name in operator dir without '`v`'. For example, updating the aqua operator from `1.0.0` to `1.0.1`

```
$ tree community-operators/aqua/ -d
community-operators/aqua/
├── 0.0.1
├── 0.0.2
├── 1.0.0
├── 1.0.1
```

## Minor (cosmetics) changes

There is some case when only some minor changes to the existing operator are needed (like a description update or an update of an icon). In this case, the pipeline will set a corresponding label and automatically handle such case.

### Allowed changes

- Only changes in an csv (*.clusterserviceversion.yaml) are allowed
- List of allowed tag changes in csv
    - `spec.description`
    - `spec.DisplayName`
    - `spec.icon`
    - `metadata.annotations.description`
    - `metadata.annotations.alm-examples`

## Operator versioning strategy 

Sometimes it is needed to change how operator versions are built to the index. This can be controlled by `ci.yaml` file. [More info](./operator-ci-yaml.md#reviewers)

## Reviewers update

While an operator is envolving over time, sometimes it is needed to change reviewers. This can be controlled by `ci.yaml` file. [More info](./operator-ci-yaml.md#operator-versioning)





