#!/usr/bin/env bash
set +x

log "################################INIT VAULT#################################"
log "Start watch loop and wait for the vault-0 pod to be running with 1/1 status in order to init properly."
max_loops=100
current_loop=0
while true
do
    check_vault_zero_is_running=$(kubectl get pods -n vault vault-0 | grep '1/1')
    log "Vault running status: ${check_vault_zero_is_running}"
    if [ -z "${check_vault_zero_is_running}" ]; then
        log "Vault Pod is not in 1 of 1 state: ${check_vault_zero_is_running}"
    else
        log "pod running lets try to init until success."
        sleep 20
        max_loops=100
        current_loop=0
        log "trying to vault init."
        kubectl exec -it vault-0 -n vault -- vault operator init > vault_output
        while [ $? -ne 0 ]; do
            sleep 10
            current_loop=$((current_loop+1))
            if [ "${current_loop}" == "${max_loops}" ]; then
                log "hit max loop count."
                exit 1
            fi
            kubectl exec -it vault-0 -n vault -- vault operator init > vault_output
        done
        break
    fi
    current_loop=$((current_loop+1))
    if [ "${current_loop}" == "${max_loops}" ]; then
        log "hit max loop count."
        exit 1
    fi
    log "${current_loop} of ${max_loops}... sleeping for 5 seconds."
    sleep 5
done

log "################################LOGIN VAULT#################################"
log "Initalize Vault, This will Produce a Token output. We will save this to variable and use it to authenticate the vault pod."

vault_output=$(cat vault_output)
VAULT_TOKEN=`echo -n $vault_output | cut -d ':' -f 7 | cut -d ' ' -f 2 | tr -d '\n'`
VAULT_TOKEN=`tr -dc '[[:print:]]' <<< "$VAULT_TOKEN"`
VAULT_TOKEN=`echo -n $VAULT_TOKEN | rev | cut -c4- | rev`
max_loops=100
current_loop=0
kubectl exec -n vault -it vault-0 -- vault login $VAULT_TOKEN
while [ $? -ne 0 ]; do
    sleep 10
    log "Trying to login to vault with token from file vault_output"
    vault_output=$(cat vault_output)
    VAULT_TOKEN=`echo -n $vault_output | cut -d ':' -f 7 | cut -d ' ' -f 2 | tr -d '\n'`
    VAULT_TOKEN=`tr -dc '[[:print:]]' <<< "$VAULT_TOKEN"`
    VAULT_TOKEN=`echo -n $VAULT_TOKEN | rev | cut -c4- | rev`
    current_loop=$((current_loop+1))
    if [ "${current_loop}" == "${max_loops}" ]; then
        log "hit max loop count."
        exit 1
    fi
    kubectl exec -n vault -it vault-0 -- vault login $VAULT_TOKEN
done
log "#####################################################"

log "################################Enable Enterprise Vault License#################################"
kubectl exec -n vault  vault-0 -- vault write sys/license text=@/opt/cloudhsm/etc/sifchain-vault-ent-license.hcli
log "#####################################################"

log "################################Enable KV-V2 Secrets Engine#################################"
log "enable kv-v2 secrets engine."
kubectl exec -n vault  vault-0 -- vault secrets enable kv-v2
log "#####################################################"

log "################################BACKUP TOKENS IN ENCRYPTED SECURE S3 Bucket vault-key-backup-$TF_VAR_profile#################################"
log "Upload vault init token backup so s3 that only the vault automation user has access too."
aws s3 cp ./vault_output s3://vault-key-backup-$TF_VAR_app_env/vault-master-keys.backup --region $TF_VAR_aws_region
log "#####################################################"

log "################################TEST VAULT IS WORKING: TEST TEST TEST TEST#################################"
kubectl exec -n vault -it vault-0 -- vault kv put kv-v2/test test="test-secret"
sleep 15
check_vault_test_secret_exist=$(kubectl exec -n vault -it vault-0 -- vault kv get kv-v2/test | grep 'test-secret')
if [ -z "${check_vault_test_secret_exist}" ]; then
  log "Your test secret was a success, vault is functioning properly."
else
  log "Your test secret didn't seem to show up vault doesn't seem to be functioning."
  exit 1
fi
