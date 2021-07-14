#!/usr/bin/env bash

function header (){
    echo "======================"
    echo "${1}"
    echo "======================"
}

function footer (){
    echo "======================"
    echo ""
}

read -p "What AWS profile do you want to use: \n $(cat ~/.aws/config | grep "\[profile" | grep -v "#") [devnet]:" AWS_PROFILE
AWS_PROFILE=${AWS_PROFILE:-devnet}


read -p "What AWS region is your Kubernetes cluster running in: [us-east-1]:" AWS_REGION
AWS_REGION=${AWS_REGION:-us-east-1}


read -p "What AWS region is your Kubernetes Cluster Name in: [sifchain-aws-devnet-us]:" AWS_CLUSTER_NAME
AWS_CLUSTER_NAME=${AWS_CLUSTER_NAME:-sifchain-aws-devnet-us}


header "Pull Kubernetes Config for ${AWS_CLUSTER_NAME} in region ${AWS_REGION} using AWS profile ${AWS_PROFILE}"
aws eks update-kubeconfig --name ${AWS_CLUSTER_NAME} --region ${AWS_REGION} --profile devnet
footer


header "Test Kubernetes Configuration"
check_kubernetes_kube_condig=$(kubectl get pods --all-namespaces | grep coredns)
if [ -z "${check_kubernetes_kube_condig}" ]; then
    echo "Kubernetes Config Doesn't Appear to be Configured Properly."
    kubectl get pods --all-namespaces
else
    echo "Kubernetes Config Working Properly."
fi
footer


header "Please tell me the values to configure vault for our governance pipeline requests and voting.."
read -p "What is the Moniker you will be setting the vault entry up for?: [devnet-us-1]:" MONIKER
MONIKER=${MONIKER:-devnet-us-1}


read -p "What is the CHAINNET you will be setting the vault entry up for?: [sifchain-devnet]:" CHAINNET
CHAINNET=${CHAINNET:-sifchain-devnet}


read -p "What is the mnemonic you will be setting the vault entry up for?: [this is where your mnemonic goes]:" MNEMONIC
MNEMONIC=${MNEMONIC:-None}
if [ "${MNEMONIC}" == "None" ]; then
    echo "YOU MUST ENTER A MNEMONIC"
    exit 1
fi
footer


header "Building the docker sifnoded image to test values import as keyring properly."
make CHAINNET=${CHAINNET} IMAGE_TAG=test BINARY=sifnoded build-image
footer


header "Run a test and import keyring with provided values"
docker run -e MNEMONIC="${MNEMONIC}" -e MONIKER="${MONIKER}" -i sifchain/sifnoded:test sh <<'EOF'
	sifnoded keys list
	yes "${MNEMONIC}" | sifnoded keys add ${MONIKER} -i --recover --keyring-backend test
exit
EOF
footer


read -p "What is the 2 character region us,eu,au,sg: [us]:" APP_REGION
APP_REGION=${APP_REGION:-us}


header "Setup Vault Entry"
kubectl exec -n vault -it vault-0 -- vault kv put kv-v2/${APP_REGION}/governance-test \
  mnemonic="${MNEMONIC}" \
  moniker="${MONIKER}"
footer

sleep 2

header "Check Vault Entry"
kubectl exec -n vault -it vault-0 -- vault kv get kv-v2/${APP_REGION}/governance-test
footer