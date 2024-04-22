# About the Sail Operator

The Sail Operator is based on the open source Istio project. The Sail implementation pulls its code from the upstream Istio master repository with no changes to the codebase.

# Prerequisites

You have deployed a cluster on OpenShift Container Platform 4.13 or later. 

You are logged in to the OpenShift Container Platform web console as a user with the `cluster-admin` role. 

# Installing the Sail Operator

1. Navigate to the operator hub. 

2. Click **Operator** -> **Operator Hub**.

3. Search on Sail. 

4. Locate the Sail Operator, and click to select it.

5. When the prompt that discusses the community operator appears, click **Continue**.

6. Verify the Sail Operator is version 3.0, and click **Install**. 

7. Use the default installation settings presented, and click **Install** to continue.

8. Click **Operators** -> **Installed Operators** to verify that the Sail Operator is installed. `Succeeded` should appear in the **Status** column. 

# Deploying the control plane

1. Create the project where Istio is going to be deployed:  

``` 
$ oc new-project istio-system
```
2. Provide privileges to the project. 

```
$ oc adm policy add-scc-to-group anyuid system:serviceaccounts:istio-system
$ oc adm policy add-scc-to-group privileged system:serviceaccounts:istio-system
```
3. In the OpenShift web console, select `istio-system` in the **Project** drop-down menu. 

4. Click the Sail Operator.

5. Click **Istio Helm Install**.

6. Click **Create IstioHelmInstall**.

7. Accept the defaults and click **Create**. This creates the control plane.

8. Click **Workloads** -> **Pods**. Verify that the pods were created. `Running` should appear in the **Status** column. If the pods were successfully created, then Istio is installed and ready for use. For more information, see the [Istio documentation](https://istio.io/latest/docs/setup/platform-setup/openshift/).

# Customizing Istio configuration

The `values` field of the `IstioHelmInstall` custom resource definition created above can be used to customize Istio configuration using Istio's `Helm` configuration values. When creating this resource with the OpenShift Container Platform web console (as described above), it will be pre-populated with configuration to enable Istio to run on OpenShift. 

To view or modify the `IstioHelmInstall` resource from the OpenShift Container Platform web console:
1. Click **Operators** -> **Installed Operators**. Under **Provided APIs**, click **Istio Helm Install**.
2. Under **Name**, click on the instance of `IstioHelmInstall` ("istiohelminstall-sample" by default).
3. Click **YAML** to view the `IstioHelmInstall` configuration and make modifications.

For a list of available configuration for the `values` field, refer to [Istio's artifacthub chart documentation](https://artifacthub.io/packages/search?org=istio&sort=relevance&page=1) for:
- [Base parameters](https://artifacthub.io/packages/helm/istio-official/base?modal=values)
- [Istiod parameters](https://artifacthub.io/packages/helm/istio-official/istiod?modal=values)
- [Gateway parameters](https://artifacthub.io/packages/helm/istio-official/gateway?modal=values)
- [CNI parameters](https://artifacthub.io/packages/helm/istio-official/cni?modal=values)
- [ZTunnel parameters](https://artifacthub.io/packages/helm/istio-official/ztunnel?modal=values)
