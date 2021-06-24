#!/usr/bin/env bash
set +x

log "Setup Variables"
log "####################################CONFIGURE VARIABLES####################################"

read -p "What is your AWS Access Key? [AKXXXX]:" AUTOMATION_USER_AWS_ACCESS_KEY
AUTOMATION_USER_AWS_ACCESS_KEY=${AUTOMATION_USER_AWS_ACCESS_KEY:-AKIAXXXX}
echo " "

read -p "What is your AWS Secret Key? [fKBa2QXXXX]:" AUTOMATION_USER_AWS_SECRET_KEY
AUTOMATION_USER_AWS_SECRET_KEY=${AUTOMATION_USER_AWS_SECRET_KEY:-fKBa2XXXXX}
echo " "

read -p "What is your AWS region? [us-east-2]:" TF_VAR_aws_region
export TF_VAR_aws_region=${TF_VAR_aws_region:-us-east-2}
echo " "

log "setup aws credential file entry."
check_credential_file=$(cat ~/.aws/credentials | grep 'sifchain-automation')
if [ -z "${check_credential_file}" ]; then
    log "configuring ~/.aws/credentials"
    echo '' >> ~/.aws/credentials
    echo '[sifchain-automation]' >> ~/.aws/credentials
    echo "aws_access_key_id = ${AUTOMATION_USER_AWS_ACCESS_KEY}" >> ~/.aws/credentials
    echo "aws_secret_access_key = ${AUTOMATION_USER_AWS_SECRET_KEY}" >> ~/.aws/credentials
    #echo "region = ${TF_VAR_aws_region}" >> ~/.aws/credentials
else
    log "credential file already configured."
fi
export TF_VAR_cred_profile=sifchain-automation
export TF_VAR_profile=devnet

read -p "What is your Kubernetes Cluster Name? [sifchain-aws-tempnet-us]:" TF_VAR_aws_kubernetes_cluster_name
export TF_VAR_aws_kubernetes_cluster_name=${TF_VAR_aws_kubernetes_cluster_name:-sifchain-aws-tempnet-us}
echo " "

read -p "What is your App Environment? [tempnet]:" TF_VAR_app_env
export TF_VAR_app_env=${TF_VAR_app_env:-tempnet}
echo " "

export TF_VAR_hsm_name=sifchain-${TF_VAR_app_env}-hsm
log "$TF_VAR_hsm_name"

read -p "Vault enterprise License String [02MV4UU43BK]:" hsm_license_key
hsm_license_key=${hsm_license_key:-02MV4UU43B}
echo " "

export TF_VAR_cluster_subnet_id=$(aws eks describe-cluster --name $TF_VAR_aws_kubernetes_cluster_name --region $TF_VAR_aws_region | jq '.cluster.resourcesVpcConfig.subnetIds' | jq -c '.[1]' | tr -d '"')
log "EKS Cluster Subnet-ID: $TF_VAR_cluster_subnet_id"

export TF_VAR_cluster_security_group_id=$(aws eks describe-cluster --name $TF_VAR_aws_kubernetes_cluster_name --region $TF_VAR_aws_region | jq '.cluster.resourcesVpcConfig.securityGroupIds' | jq -c '.[]' | tr -d '"')
log "EKS Cluster SecurityGroup-ID: $TF_VAR_cluster_security_group_id"

export TF_VAR_cluster_node_security_group_id=$(aws ec2 describe-security-groups --region $TF_VAR_aws_region --filters Name=tag:Name,Values=$TF_VAR_aws_kubernetes_cluster_name-eks_cluster_sg | jq '.SecurityGroups' | jq -c '.[]' | jq '.GroupId' | tr -d '"' | tr -d '\n')
log "EKS Cluster SecurityGroup-ID: $TF_VAR_cluster_node_security_group_id"

export TF_VAR_cluster_level_security_group_id=$(aws ec2 describe-security-groups --region $TF_VAR_aws_region --filters Name=tag:"aws:eks:cluster-name",Values=$TF_VAR_aws_kubernetes_cluster_name | jq '.SecurityGroups' | jq -c '.[]' | jq '.GroupId' | tr -d '"' | tr -d '\n')
log "EKS Node SecurityGroup-ID: $TF_VAR_cluster_level_security_group_id"

HSM_USER=vault_user
log "HSM Vault User: $HSM_USER"

read -p "HSM Password for Vaults Cloud HSM User? [GHScnsdj523jfai48cjaFJ]:" HSM_PASSWORD
HSM_PASSWORD=${HSM_PASSWORD:-GHScnsdj523jfai48cjaFJ}
echo " "

read -p "HSM Admin Password for Vaults Cloud HSM User? [nsSDJRwe832Wnfwmle34Dks]:" HSM_ADMIN_PASSWORD
HSM_ADMIN_PASSWORD=${HSM_ADMIN_PASSWORD:-nsSDJRwe832Wnfwmle34Dks}
echo " "

read -p "Docker Hub Username for Kubernetes Docker Secret for Vault? [gzukel]:" DOCKER_USERNAME
DOCKER_USERNAME=${DOCKER_USERNAME:-gzukel}
echo " "

read -p "Docker Hub Generated Access Token for Kubernetes Docker Secret for Vault? [930d1007-XXXXX]:" DOCKER_PASSWORD
DOCKER_PASSWORD=${DOCKER_PASSWORD:-c354af65-XXXXX}
echo " "

DOCKER_SERVER=https://index.docker.io/v2/
log "Docker Server: ${DOCKER_SERVER}"

read -p "Docker Hub Email Associated with Access Token? [grant@sifchain.finance]:" DOCKER_EMAIL
DOCKER_EMAIL=${DOCKER_EMAIL:-grant@sifchain.finance}
echo " "

read -p "Is IS_DEBUG? [true]:" IS_DEBUG
IS_DEBUG=${IS_DEBUG:-true}
echo " "