# Get started

Our Terraform modules are adaptable for different cloud providers from where the Kubernetes cluster is provisioned. This deployment will initialize a cluster environment prepared for deploying Sifnode.

# Requirements

- AWS account
- CLI and AWS credentials configured
- kubectl

# Provision and configure

Infrastructure provision with the `make` command.

## AWS

During the deploy process you will be asked to enter `cluster name` and `AWS region`.

```
$ make aws
```

Once the EKS cluster is provisioned, you need to setup your local config context for kubectl. Update the following command with your cluster name and region.

```
$ aws eks --region <cluster_region> update-kubeconfig --name <cluster_name>
$ kubectl version
```

To verify the status of your cluster run:

```
$ kubectl get nodes
```

# Destroy

Destroy the whole infrastructure with the following command:

## AWS

```
$ make aws-destroy
```
