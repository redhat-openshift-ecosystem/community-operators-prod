# About the Sail Operator

The Sail Operator is based on the open source Istio project. The Sail implementation pulls its code from the upstream Istio main repository with no changes to the codebase.

## Prerequisites

You have deployed a cluster on OpenShift Container Platform 4.13 or later.

You are logged in to the OpenShift Container Platform web console as a user with the `cluster-admin` role.

You have access to the OpenShift CLI (oc).

## Installing the Sail Operator

1. Navigate to the OperatorHub.

2. Click **Operator** -> **Operator Hub**.

3. Search on Sail.

4. Locate the Sail Operator, and click to select it.

5. When the prompt that discusses the community operator appears, click **Continue**.

6. Verify the Sail Operator is version 3.0, and click **Install**.

7. Use the default installation settings presented, and click **Install** to continue.

8. Click **Operators** -> **Installed Operators** to verify that the Sail Operator is installed. `Succeeded` should appear in the **Status** column.

## Deploying Istio

1. Create the project where Istio is going to be deployed:  

    ```sh
    $ oc new-project istio-system
    ```

1. In the OpenShift web console, select `istio-system` in the **Project** drop-down menu.

1. Click the Sail Operator.

1. Click **Istio**.

1. Click **Create Istio**.

1. Accept the defaults and click **Create**. This creates the Istio control plane.

1. Click **Workloads** -> **Pods**. Verify that the pods were created. `Running` should appear in the **Status** column. If the pods were successfully created, then Istio is installed and ready for use. For more information, see the upstream [Istio documentation](https://istio.io/latest/docs/setup/platform-setup/openshift/).

## Customizing Istio configuration

The `values` field of the `Istio` custom resource definition, which was created when the control plane was deployed, can be used to customize Istio configuration using Istio's `Helm` configuration values. When you create this resource using the OpenShift Container Platform web console, it is pre-populated with configuration settings to enable Istio to run on OpenShift.

To view or modify the `Istio` resource from the OpenShift Container Platform web console:

1. Click **Operators** -> **Installed Operators**.
1. Click **Istio** in the **Provided APIs** column.
1. Click `Istio` instance, "istio-sample" by default, in the **Name** column.
1. Click **YAML** to view the `Istio` configuration and make modifications.

For a list of available configuration for the `values` field, refer to [Istio's artifacthub chart documentation](https://artifacthub.io/packages/search?org=istio&sort=relevance&page=1) for:

- [Base parameters](https://artifacthub.io/packages/helm/istio-official/base?modal=values)
- [Istiod parameters](https://artifacthub.io/packages/helm/istio-official/istiod?modal=values)
- [Gateway parameters](https://artifacthub.io/packages/helm/istio-official/gateway?modal=values)
- [CNI parameters](https://artifacthub.io/packages/helm/istio-official/cni?modal=values)
- [ZTunnel parameters](https://artifacthub.io/packages/helm/istio-official/ztunnel?modal=values)

## Installing the Bookinfo Application

You can use the `bookinfo` example application to explore service mesh features. Using the `bookinfo` application, you can easily confirm that requests from a web browser pass through the mesh and reach the application.

The `bookinfo` application displays information about a book, similar to a single catalog entry of an online book store. The application displays a page that describes the book, lists book details (ISBN, number of pages, and other information), and book reviews.

The `bookinfo` application is exposed through the mesh, and the mesh configuration determines how the microservices comprising the application are used to serve requests. The review information comes from one of three services: `reviews-v1`, `reviews-v2` or `reviews-v3`. If you deploy the `bookinfo` application without defining the `reviews` virtual service, then the mesh uses a round robin rule to route requests to a service.

By deploying the `reviews` virtual service, you can specify a different behavior. For example, you can specify that if a user logs into the `bookinfo` application, then the mesh routes requests to the `reviews-v2` service, and the application displays reviews with black stars. If a user does not log into the `bookinfo` application, then the mesh routes requests to the `reviews-v3` service, and the application displays reviews with red stars.

For more information, see [Bookinfo Application](https://istio.io/latest/docs/examples/bookinfo/) in the upstream Istio documentation.

## Gateway Configuration

The Sail Operator does not deploy Ingress or Egress Gateways. Gateways are not part of the control plane. As a security best-practice, Ingress and Egress Gateways should be deployed in a different namespace than the namespace that contains the control plane.

You can deploy gateways using either the Gateway API or Gateway Injection methods. Both are well documented in the Istio documentation.

- To use Gateway API, follow the instructions in the [Getting Started with Istio and Kubernetes Gateway API](https://preliminary.istio.io/latest/docs/setup/additional-setup/getting-started/) page.
- To use Gateway Injection, use the `Helm` method described in the [Installing Gateways](https://preliminary.istio.io/latest/docs/setup/additional-setup/gateway/#deploying-a-gateway) page.
