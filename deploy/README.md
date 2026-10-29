# Cluster deployment

Use `python3 deploy/cluster.py --help` to scaffold and manage AWS EKS or Azure AKS,
then deploy the shared sifnode Helm chart. Each environment keeps its own Terraform
state and kubeconfig under the ignored `.live/` directory.

See [Running sifnode on Azure AKS](../docs/tutorials/running-sifnode-on-aks.md)
for authentication, cluster provisioning, genesis setup, node deployment, AWS
selection, teardown and local verification.
