#!/usr/bin/env bash
set +x
read -p "You are about to run the terraform template to deploy the AWS Cloud HSM. This will configure the HSM and the Cluster. The automation will also create a signed cert from the HSM cluster cert. Are you sure you want to do this? [y/n]:" terraform_yes_no
terraform_yes_no=${terraform_yes_no:y}
if [ "${terraform_yes_no}" == "y" ]; then

    if [ "$IS_DEBUG" == "true" ]; then
        log ""
        log "Output main.tf for review"
        cat main.tf
    fi

    log ""
    log "Terraform Init."
    terraform init .

    log ""
    log "Apply Terraform File."
    terraform apply .
fi