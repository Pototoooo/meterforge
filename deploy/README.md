# Deploy MeterForge

## Deploy MeterForge to a local Kubernetes cluster

## Prerequisites

- [docker](https://www.docker.com/)
- [kind](https://kind.sigs.k8s.io/)
- [helm](https://helm.sh/)

## 1. Check out this repository

```shell
git clone git@github.com:Pototoooo/meterforge.git
cd meterforge/deploy
```

## 2. Setup local cluster

```shell
kind create cluster --config ./kind.yaml
```

## 3. Install MeterForge via Helm

```shell
helm upgrade --install --dependency-update -f ./charts/meterforge/values.example.yaml meterforge ./charts/meterforge
```

Once the `meterforge-api` pod is ready, we can use `port-forward` to access the API.
This might take a few minutes.

```shell
kubectl port-forward svc/meterforge-api 8888:80
```
