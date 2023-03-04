# Kubernetes Chaos Monkey

A simple program that when running in Kubernetes cluster deletes a random Pod in a given namespace on a given period of time.

## Program structure

This is a simple Go program that uses the [Go Kubernetes Client](https://github.com/kubernetes/client-go) to communicate with the cluster in which the app is running,
read the pods at a given namespace and delete one of them at random.
It then sleep for a configurable period and it does the same thing again.
It runs forever.

## Configuration

The program has 2 configurable fields:

- `K8S_CHAOS_NAMESPACE` - the namespace in which it will delete pods.
  Default: `default`.
- `K8S_CHAOS_PERIOD_SECONDS` - the period (in seconds) which it will wait after it deletes a Pod.
  Default: `10`.

## Container image

In order to run this in Kubernetes we need to build a container image out of it.
There is a [`Dockerfile`](./Dockerfile) in the root directory which we use for that purpose.
It builds a container image based on the `golang:1.19-alpine` base image.

If you want to build your own image you can do it with this command:

```shell
docker build -t <ACCOUNT>/<REPO>:<TAG> -f Dockerfile .
```

Or you can use the one which I already built and pushed into my public Docker Hub profile: [asankov/k8s-chaos-monkey](https://hub.docker.com/r/asankov/k8s-chaos-monkey/tags).

**NOTE:** Using `latest` is a bad practice, because when we use `latest` tag we don't have any reproducability over what runs in our cluster.
If you want to pull this image use the `0.1` tag.

## Kubernetes

If you want to run this in Kubernetes follow the instructions in [this folder](./k8s/).

## Improvements

Future improvements for this app could be:

- using multi-stage container build with [`scratch`](https://hub.docker.com/_/scratch) or [`distroless`](https://github.com/GoogleContainerTools/distroless) image.
  Since this container contains just a Go binary that is static and does not need any other OS dependencies, we can use scratch or distroless containers
  that are much more light-weight that the alpine one we are currently using.
- more configuration options - for example, we can exclude pods, based on labels/annonations.
- using a logging library, that will give us the ability to customize the verbosity of the logs
