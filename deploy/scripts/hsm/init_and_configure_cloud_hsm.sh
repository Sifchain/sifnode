#!/usr/bin/env bash
set +x

log "################################Wait for HSM to be Initialized#################################"
log "Wait for the cluster to be in initalized state.."
max_loops=100
current_loop=0
hsm_cluster_id=$(cat cluster_id)
log "HSM Cluster ID ${hsm_cluster_id}"
while true
do
    if [ "${IS_DEBUG}" == "true" ]; then
      aws cloudhsmv2 describe-clusters --region $TF_VAR_aws_region --filters clusterIds=$hsm_cluster_id
    fi
    check_init=$(aws cloudhsmv2 describe-clusters --region $TF_VAR_aws_region --filters clusterIds=$hsm_cluster_id | grep -E -- 'INITIALIZED|ACTIVE')
    log "Is Cluster Initialized: ${check_init}"
    if [ -z "${check_init}" ]; then
        current_loop=$((current_loop+1))
        if [ "${current_loop}" == "${max_loops}" ]; then
            log "hit max loop count."
            exit 1
        fi
        log "${current_loop} of ${max_loops}... sleeping for 5 seconds."
        sleep 10
    else
        check_is_not_initalizing=$(aws cloudhsmv2 describe-clusters --region $TF_VAR_aws_region --filters clusterIds=$hsm_cluster_id | grep 'INITIALIZE_IN_PROGRESS')
        if [ -z "${check_is_not_initalizing}" ]; then
            log "Cluster initalized. Continuing."
            break
        fi
        log "Still initializing"
        current_loop=$((current_loop+1))
        if [ "${current_loop}" == "${max_loops}" ]; then
            log "hit max loop count."
            exit 1
        fi
        log "${current_loop} of ${max_loops}... sleeping for 5 seconds."
        sleep 10
    fi
done

log "Get the CLOUD HSM IP for your newly created cloudHSM cluster."
HSM_IP=$(aws cloudhsmv2 describe-clusters --region $TF_VAR_aws_region --filters clusterIds=$hsm_cluster_id | jq '.Clusters' | jq -c '.[0]' | jq -c '.Hsms' | jq -c '.[0]' | jq '.EniIp' | tr -d '"')
if [ "$IS_DEBUG" == "true" ]; then
    log "CloudHSM IP: ${HSM_IP}"
fi

if [ -z "HSM_IP" ]; then
    echo "Cloud HSM IP ${HSM_IP}"
    exit 1
fi

log "Get the contents of the generated CA certificate for the CloudHSM.."
CA_CONTENTS=$(cat customerCA.crt | base64 | tr -d '\n')

log " "
log " "
read -p "You are about to enter the configure step for the HSM automation. This is where you will have to enter some manual commands to configure the HSM users. Do you want to proceed? [y/n]:" configure_yes_no
configure_yes_no=${configure_yes_no:y}
log " "
if [ "${configure_yes_no}" == "y" ]; then
    log "#######################################"
    log " "
    log "#####################USE KUBERENETES TEMP SHELL TO CONFIGURE HSM USERS#####################"
    log "Open temporary kubernetes shell to configure the cloud hsm users. This is the only manual part please pay attention to the directed copy and paste."
    kubectl run test-shell --rm -i --tty --image ubuntu:bionic-20210512 -- bash -c "\
        apt-get update -y && \
        echo 'ca contents' && \
        echo $CA_CONTENTS && \
        apt-get install nano wget unzip ca-certificates gnupg openssl libpcap-dev dumb-init tzdata -y && \
        wget https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh && \
        bash Miniconda3-latest-Linux-x86_64.sh -ab && \
        export PATH=~/miniconda3/bin:${PATH} && \
        wget https://s3.amazonaws.com/cloudhsmv2-software/CloudHsmClient/Bionic/cloudhsm-client_latest_u18.04_amd64.deb && \
        wget https://s3.amazonaws.com/cloudhsmv2-software/CloudHsmClient/Bionic/cloudhsm-client-pkcs11_latest_u18.04_amd64.deb && \
        apt install ./cloudhsm-client_latest_u18.04_amd64.deb -y && \
        apt install ./cloudhsm-client-pkcs11_latest_u18.04_amd64.deb -y && \
        echo \"${CA_CONTENTS}\" | base64 --decode > /opt/cloudhsm/etc/customerCA.crt && \
        apt-get install opensc -y && \
        /opt/cloudhsm/bin/configure -a ${HSM_IP} && \
        echo 'You will run the following cli commands in the aws-cloudhsm cli. This will configure the users on aws.' && \
        echo 'There is no way to automate this so you must run this manually.' && \
        echo 'aws-cloudhsm> loginHSM PRECO admin password' && \
        echo 'aws-cloudhsm> changePswd PRECO admin {HSM_ADMIN_PASSWORD}' && \
        echo 'aws-cloudhsm> logoutHSM' && \
        echo 'aws-cloudhsm> loginHSM CO admin {HSM_ADMIN_PASSWORD}' && \
        echo 'aws-cloudhsm> createUser CU {HSM_USER} {HSM_PASSWORD}' && \
        echo 'aws-cloudhsm> quit' && \
        export HSM_USER=vault_user && \
        echo HSM_USER: $HSM_USER && \
        echo HSM_PASSWORD: $HSM_PASSWORD && \
        echo HSM_ADMIN_PASSWORD: $HSM_ADMIN_PASSWORD && \
        /opt/cloudhsm/bin/cloudhsm_mgmt_util /opt/cloudhsm/etc/cloudhsm_mgmt_util.cfg && \
        exit"
fi
