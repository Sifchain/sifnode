# Get started

We have a Helm Chart for deploying the Sifnode daemon on our already launched Kubernetes cluster.

# Requirements

- Running Kubernetes cluster
- Kubectl configured, ready and connected to running cluster
- Helm 3 binary (Helm 3 can also be installed with `make helm`).
- GitHub account for pulling the Docker Image from private Container Registry.
- GitHub token

# Deploy

Installing Helm dependencies before deploying the Sifnode daemon

```
$ make repos
```

You need to generate an encrypted value from your GitHub credentials and put the content in values.yaml

```
$ export GITHUB_USERNAME=<username>
$ export GITHUB_TOKEN=<token>
$ kubectl create secret docker-registry --dry-run=client ghcr --docker-server=docker.pkg.github.com --docker-username=${GITHUB_USERNAME} --docker-password=${GITHUB_TOKEN} -o jsonpath='{.data.\.dockerconfigjson}'
```

Get the output from the last command and put it in the Helm's values.yaml file as a value of `dockerconfigjson`. Once the encrypted value is added to the values.yaml file we can install our chart:

```
$ make deploy-sifnode nodename=<nodename>
```

# Destroy

For destroying the sifnode chart just execute the destroy command:

```
$ make destroy-sifnode nodename=<nodename>
```
