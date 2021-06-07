#!/usr/bin/env bash

read -p "What is your app name [sifnode]:" APP_NAME
APP_NAME=${APP_NAMES:-sifnode}

read -p "What is your app environment [devnet]:" APP_ENV
APP_ENV=${APP_ENV:-devnet}

read -p "What is your app region [us]:" APP_REGION
APP_REGION=${APP_REGION:-us}

read -p "What is your aws profile [devnet]:" AWS_PROFILE
AWS_PROFILE=${AWS_PROFILE:-devnet}

read -p "What is your aws region [us-east-1]:" AWS_REGION
AWS_REGION=${AWS_REGION:-us-east-1}

read -p "What is your kubernetes cluster name [sifchain-aws-devnet-us]:" CLUSTER_NAME
CLUSTER_NAME=${CLUSTER_NAME:-sifchain-aws-devnet-us}

read -p "What is your application image name [sifchain/sifnoded]:" IMAGE_NAME
IMAGE_NAME=${IMAGE_NAME:-sifchain/sifnoded}

read -p "What is your application image tag [testnet-genesis]:" IMAGE_TAG
IMAGE_TAG=${IMAGE_TAG:-testnet-genesis}

read -p "What is your application peer address [1b02f2eb065031426d37186efff75df268bb9097@54.164.57.141:26656]:" PEER_ADDRESS
PEER_ADDRESS=${PEER_ADDRESS:-1b02f2eb065031426d37186efff75df268bb9097@54.164.57.141:26656}

read -p "What is the directory to your helm chart [deploy/helm/sifnode-vault]:" TEMPLATE_ROOT
TEMPLATE_ROOT=${TEMPLATE_ROOT:-deploy/helm/sifnode-vault}

read -p "What is the name of the helm template values file [template.values.yaml]:" TEMPLATE_FILE
TEMPLATE_FILE=${TEMPLATE_FILE:-template.values.yaml}

read -p "What is the name of the values file you want to generate [generated.values.yaml]:" GENERATED_FILE
GENERATED_FILE=${GENERATED_FILE:-generated.values.yaml}

read -p "What is your chainnet [sifchain-devnet]:" CHAINNET
CHAINNET=${CHAINNET:-sifchain-devnet}

read -p "What is your genesisURL [https://raw.githubusercontent.com/Sifchain/networks/master/testnet/sifchain-devnet/genesis.json]:" genesisURL
genesisURL=${genesisURL:-https://raw.githubusercontent.com/Sifchain/networks/master/testnet/sifchain-devnet/genesis.json}

read -p "What is your mnemonic [your mnemonic goes here]:" mnemonic
mnemonic=${mnemonic:-your mnemonic goes here}

read -p "What is your kubernetes namespace goes here [sifnode]:" APP_NAMESPACE
APP_NAMESPACE=${APP_NAMESPACE:-sifnode}

read -p "What is your moniker goes here [devnet-us-1]:" MONIKER
MONIKER=${MONIKER:-devnet-us-1}

echo "APP_ENV=$APP_ENV"
echo "APP_REGION=$APP_REGION"
echo "APP_NAME=$APP_NAME"
echo "AWS_PROFILE=$AWS_PROFILE"
echo "AWS_REGION=$AWS_REGION"
echo "CLUSTER_NAME=$CLUSTER_NAME"
echo "IMAGE_NAME=$IMAGE_NAME"
echo "IMAGE_TAG=$IMAGE_TAG"
echo "PEER_ADDRESS=$PEER_ADDRESS"
echo "TEMPLATE_ROOT=$TEMPLATE_ROOT"
echo "TEMPLATE_FILE=$TEMPLATE_FILE"
echo "GENERATED_FILE=$GENERATED_FILE"
echo "CHAINNET=$CHAINNET"
echo "genesisURL=$genesisURL"
echo "mnemonic=$mnemonic"
echo "MONIKER=$MONIKER"
echo "APP_NAMESPACE=$APP_NAMESPACE"
export app_region=${APP_REGION}
export app_env=${APP_ENV}

echo "Pull Kubernetes Config"
echo "======================"
aws eks update-kubeconfig \
    --name ${CLUSTER_NAME} \
    --region ${AWS_REGION} \
    --profile ${AWS_PROFILE}
echo " "

echo "List Pods to ensure Kubeconfig is working"
echo "========================================="
kubectl get pods --all-namespaces
echo " "

echo "Create Vault Entry"
read -p "Do you want me to create the vault entry y/n [y]:" CREATE_VAULT
CREATE_VAULT=${CREATE_VAULT:-y}
if [ "${CREATE_VAULT}" == "y" ]; then
    echo "=================="
    kubectl exec -n vault -it vault-0 -- vault kv put kv-v2/${APP_REGION}/${APP_ENV}/${APP_NAME} \
        CHAINNET="${CHAINNET}" \
        govVotingPeriod="x" \
        govMaxDepositPeriod="x" \
        minimumGasPrices="x" \
        adminOracleAddress="x" \
        adminCLPAddresses="x" \
        bondAmount="x" \
        mintAmount="x" \
        genesisURL="${genesisURL}" \
        mnemonic="${mnemonic}" \
        MONIKER="${MONIKER}"
    echo " "
else
    exit 0
fi

echo "Check secret was created"
echo "========================"
kubectl exec -n vault -it vault-0 -- vault kv get kv-v2/${APP_REGION}/${APP_ENV}/${APP_NAME}
echo " "

echo "Copy the kubeconfig so the rake tasks can work for setting up the application."
echo "=============================================================================="
cp ~/.kube/config ./kubeconfig
echo " "

echo "Check application is setup in vault."
echo "===================================="
rake "cluster:vault:check_application_configured[${APP_ENV}, ${APP_REGION}, ${APP_NAME}]"
echo " "

echo "Create Application Vault Policy"
echo "==============================="
read -p "Do you want me to create application vault policy y/n [y]:" CREATE_APPLICATION_VAULT_POLICY
CREATE_APPLICATION_VAULT_POLICY=${CREATE_APPLICATION_VAULT_POLICY:-y}
if [ "${CREATE_APPLICATION_VAULT_POLICY}" == "y" ]; then
    echo "=================="
    rake "cluster:vault:create_vault_policy[${APP_REGION}, ${APP_NAMESPACE}, ${IMAGE_NAME}, ${IMAGE_TAG}, ${APP_ENV}, ${APP_NAME}]"
    echo " "
else
    exit 0
fi

echo "Enable Kubernetes for Vault if Not Enabled"
echo "=========================================="
read -p "Do you want me to enable vault for kubernetes if its not already y/n [y]:" ENABLE_KUBERNETES
ENABLE_KUBERNETES=${ENABLE_KUBERNETES:-y}
if [ "${ENABLE_KUBERNETES}" == "y" ]; then
    echo "=================="
    rake "cluster:vault:enable_kubernetes[]"
    echo " "
else
    exit 0
fi

echo "Configure Application in Vault"
echo "=============================="
read -p "Do you want me to configure your application for vault? y/n [y]:" CONFIGURE_APPLICATION_FOR_VAULT
CONFIGURE_APPLICATION_FOR_VAULT=${CONFIGURE_APPLICATION_FOR_VAULT:-y}
if [ "${CONFIGURE_APPLICATION_FOR_VAULT}" == "y" ]; then
    echo "=================="
    rake "cluster:vault:configure_application[${APP_NAMESPACE}, ${IMAGE_NAME}, ${IMAGE_TAG}, ${APP_ENV}, ${APP_NAME}]"
    echo " "
else
    exit 0
fi

echo "Remove Kubeconfig"
echo "================="
rm -rf ./kubeconfig
echo " "

echo "Deploy Peered Sifnode In Vault"
echo "=============================="
read -p "Do you want me to deploy sifnode? y/n [y]:" DEPLOY_SIFNODE
DEPLOY_SIFNODE=${DEPLOY_SIFNODE:-y}
if [ "${DEPLOY_SIFNODE}" == "y" ]; then
    echo "=================="
    rake "sifnode:standalone:deploy:peer_vault[${APP_NAMESPACE}, ${IMAGE_NAME}, ${IMAGE_TAG}, ${PEER_ADDRESS}, ${TEMPLATE_ROOT}/${TEMPLATE_FILE}, ${TEMPLATE_ROOT}/${GENERATED_FILE}]"
    echo " "

    echo "Remove Generated File"
    echo "================="
    rm -rf ${TEMPLATE_ROOT}/${GENERATED_FILE}
else
    exit 0
fi

