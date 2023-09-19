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
