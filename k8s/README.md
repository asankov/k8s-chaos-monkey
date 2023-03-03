# Kubernetes manifests

In this folder you can find the Kubernete manifests to run the application.
There are 2 files:

[`chaos-monkey-deployment.yaml`](./chaos-monkey-deployment.yaml) - this is the file that container the manifests to run the app.
There is a `Deployment` for the app container.
Also there is a `ClusterRole` and `ClusterRoleBinding`.
There are needed to give permissions to the Chaos Monkey Deployment to READ and DELETE pods in other namespaces.
They are cluster scoped, so they will work for EVERY namespace.

[`workloads-deployments.yaml`](./workloads-deployments.yaml) - this is a random workload deployment (in this case, `nginx` with 3 replicas).
This is needed in order for us to be able to test the Chaos Monkey app.

## Deploy

To deploy the application:

1. Point your Kubeconfig to the cluster you want to test with

2. Deploy the workloads

    ```shell
    kubectl apply -f workloads-deployments.yaml
    ```

3. Deploy the Chaos Monkey app:

    ```shell
    kubectl apply -f chaos-monkey-deployment.yaml
    ```

4. Observe the logs of the Chaos Monkey app and the state of the pods in the other namespace:

    ```shell
    $ kubectl logs deploy/k8s-chaos-monkey
    There are 3 pods in the [chaos] namespace
    Chose to delete pod [nginx-deployment-7fb96c846b-9qtxx]
    Succesfully deleted the pod [nginx-deployment-7fb96c846b-9qtxx]
    There are 3 pods in the [chaos] namespace
    Chose to delete pod [nginx-deployment-7fb96c846b-z6kpt]
    Succesfully deleted the pod [nginx-deployment-7fb96c846b-z6kpt]
    There are 3 pods in the [chaos] namespace
    Chose to delete pod [nginx-deployment-7fb96c846b-tnq9z]
    Succesfully deleted the pod [nginx-deployment-7fb96c846b-tnq9z]
    There are 3 pods in the [chaos] namespace
    Chose to delete pod [nginx-deployment-7fb96c846b-qhpkf]
    Succesfully deleted the pod [nginx-deployment-7fb96c846b-qhpkf]
    There are 3 pods in the [chaos] namespace
    ```

    ```shell
    kubectl get pods -n chaos -w
    NAME                                READY   STATUS    RESTARTS   AGE
    nginx-deployment-7fb96c846b-bnzpx   1/1     Running   0          4s
    nginx-deployment-7fb96c846b-w6fn7   1/1     Running   0          64s
    nginx-deployment-7fb96c846b-zg64f   1/1     Running   0          14s
    nginx-deployment-7fb96c846b-zg64f   1/1     Terminating   0          20s
    nginx-deployment-7fb96c846b-lz99z   0/1     Pending       0          0s
    nginx-deployment-7fb96c846b-lz99z   0/1     Pending       0          0s
    nginx-deployment-7fb96c846b-lz99z   0/1     ContainerCreating   0          0s
    nginx-deployment-7fb96c846b-zg64f   0/1     Terminating         0          21s
    nginx-deployment-7fb96c846b-zg64f   0/1     Terminating         0          21s
    nginx-deployment-7fb96c846b-zg64f   0/1     Terminating         0          21s
    nginx-deployment-7fb96c846b-lz99z   1/1     Running             0          1s
    ```
