#!/usr/bin/env bash
set +x

cd hsm/

source helpers.sh

source get_inputs.sh

source setup_aws_s3_bucket.sh

aws eks update-kubeconfig --name ${TF_VAR_aws_kubernetes_cluster_name} --region $TF_VAR_aws_region

source build_aws_cloud_hsm_with_terraform.sh

source init_and_configure_cloud_hsm.sh

log "Configure the Enterprise Vault HSM Values File"

sed -e "s/\-\=HSM_USER\=\-/${HSM_USER}/g" -e "s/\-\=HSM_PASSWORD\=\-/${HSM_PASSWORD}/g" template.vault.values.yml > hsm_values.yml
if [ "${IS_DEBUG}" == "true" ]; then
    cat hsm_values.yml
fi

log "Configure the Vault Kubernetes secret"
eval "echo \"$(< template.vault.secret.yml)\"" > hsm_secret.yml
if [ "${IS_DEBUG}" == "true" ]; then
    cat hsm_secret.yml
fi

read -p "You are about to create the Kubernetes Namespace Vault do you want to continue? [y/n]:" kube_yes_no
kube_yes_no=${kube_yes_no:y}
if [ "${kube_yes_no}" == "y" ]; then
    log "Create Kubernetes Namespace for Vault."
    kubectl create namespace vault
fi
kube_yes_no=""

read -p "You are about to create the Kubernetes docker secret? [y/n]:" kube_yes_no
kube_yes_no=${kube_yes_no:y}
if [ "${kube_yes_no}" == "y" ]; then
    log "Create Kubernetes Docker Pull Secret."
    kubectl create secret docker-registry vault-docker-secret \
        --docker-server="https://index.docker.io/v2/" \
        --docker-username="${DOCKER_USERNAME}" \
        --docker-password="${DOCKER_PASSWORD}" \
        --docker-email="${DOCKER_EMAIL}" -n vault
fi
kube_yes_no=""

read -p "You are about to create the Kubernetes Vault secret? [y/n]:" kube_yes_no
kube_yes_no=${kube_yes_no:y}
if [ "${kube_yes_no}" == "y" ]; then
    log "Apply the Kubernetes HSM Secret."
    kubectl apply -f hsm_secret.yml -n vault
fi
kube_yes_no=""

read -p "You are about to add Hashicorp helm repo and deploy enterprise vault via helm do you wish to continue? [y/n]:" kube_yes_no
kube_yes_no=${kube_yes_no:y}
if [ "${kube_yes_no}" == "y" ]; then
    log "Add the Hashicorp Helm Repo if it doesn't exist."
    helm repo add hashicorp https://helm.releases.hashicorp.com

    log "HELM unstall our custom ent vault deployment with hsm_vaules.yml"
    helm upgrade --install vault hashicorp/vault \
        --namespace vault \
        -f hsm_values.yml \
        --set image.repository=sifchain/vault \
        --set image.tag=1.2.6

    if [ "$IS_DEBUG" == "true" ]; then
        log "Output Kubernetes Logs for Vault 0 Pod"
        kubectl logs -n vault vault-0
    fi
fi
kube_yes_no=""

log "Rest for 30 seconds just to give things some time to start up. This Gives kubernetes time to schedule the new pod."
sleep 30

source init_and_configure_vault.sh