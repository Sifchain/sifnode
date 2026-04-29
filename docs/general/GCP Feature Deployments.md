# Sifnode Feature Deployments on GCP

## Push

Each feature branch will deploy a single Sifnode node to a testnet GKE cluster on Google Cloud, with each push.

### Requirements

If you wish to interact with your branch's namespace, you will need to install the following:

* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl)
* [gcloud SDK](https://cloud.google.com/sdk/docs/install)

### Namespaces

Each feature branch will have its own [namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces).

The namespace is derived from the branch name, however all slashes (`/`) will be converted to hyphens (`-`) and it will be converted to lowercase.

For example: `feature/Some-Branch` will have a namespace of `feature-some-branch`.

### kubectl

In order to interact with the cluster and your namespace, you'll need to ensure that the dependencies above have been installed.

Because the node endpoint is not exposed as part of the Github actions workflow, there are a few commands you'll need to run to obtain this. You'll require this when wanting to interact with your newly deployed node.

1. Configure gcloud:

```
gcloud init
```

*Follow the on-screen instructions to authenticate and configure the SDK.*

2. Create a new `kubeconfig`:

```
gcloud container clusters get-credentials <cluster-name> --zone <zone> --project <project-id>
```

3. Get your node's IP/s:

```
kubectl get svc -n <namespace>
```

where `<namespace>` is your namespace name, as above.

The command will display those running services in the provided namespace:

```
NAME         TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
sifnoded     LoadBalancer   10.0.0.1       35.232.1.1    1317:30336/TCP   20h
sifnoded     LoadBalancer   10.0.0.2       35.232.1.2    26656-26657:30840-30841/TCP   20h
```

You can then use the `EXTERNAL-IP` address to interact with the services, on the standard RPC/P2P (26656/26657) and rest (1317) ports.

#### Troubleshooting

If you're having trouble with your namespace, you can always query all available namespaces on the cluster:

```
kubectl get svc --all-namespaces
```

## Merging into Develop

When your feature branch is merged, the namespace will be destroyed.
