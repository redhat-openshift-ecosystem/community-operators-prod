# Where to contribute

Once you have forked the upstream repo, you will require to add your Operator Bundle to the forked repo. The forked repo will have directory structure similar to the structure outlined below.

```bash
├── config.yaml
├── operators
│   └── new-operator
│       ├── 0.0.102
│       │   ├── manifests
│       │   │   ├── new-operator.clusterserviceversion.yaml
│       │   │   ├── new-operator-controller-manager-metrics-service_v1_service.yaml
│       │   │   ├── new-operator-manager-config_v1_configmap.yaml
│       │   │   ├── new-operator-metrics-reader_rbac.authorization.k8s.io_v1_clusterrole.yaml
│       │   │   └── tools.opdev.io_demoresources.yaml
│       │   ├── metadata
│       │   │   └── annotations.yaml
│       │   └── tests
│       │       └── scorecard
│       │           └── config.yaml
│       └── ci.yaml
└── README.md
```

Follow the `operators` directory in the forked repo. Add your Operator Bundle under this `operators` directory following the example format.
1. Under the `operators` directory, create a new directory with the name of your operator.
1. Inside of this newly created directory add your `ci.yaml`.
1. Also, under the new directory create a subdirectory for each version of your Operator.
1. In each version directory there should be a `manifests/` directory containing your OpenShift yaml files, a `metadata/` directory containing your `annotations.yaml` file, and a `tests/` directory containing the required `config.yaml` file for the preflight tests.

>**Note** To learn more about preflight tests please follow this [link](https://github.com/redhat-openshift-ecosystem/openshift-preflight?tab=readme-ov-file#preflight).

For partners and ISVs, certified operators can now be submitted via connect.redhat.com. If you have submitted your Operator there already, please ensure your submission here uses a different package name (refer to the README for more details).
