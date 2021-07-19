#!/usr/bin/env bash
set +x
log "#############Create AWS Bucket############"
log "Check if s3 bucket vault-key-backup-${TF_VAR_app_env}  exists."

bucket_check=$(aws s3api list-buckets --region $TF_VAR_aws_region | grep vault-key-backup-$TF_VAR_app_env )
if [ "$IS_DEBUG" == "true" ]; then
    log "key vault backup s3 bucket exist check ${bucket_check}"
fi
if [ -z "${bucket_check}" ]; then
    read -p "You are about to create the s3 bucket vault-key-backup-$TF_VAR_app_env. This bucket will be utilized to backup your vault master keys are you sure you wan to this? [y/n]:" s_three_yes_no
    export s_three_yes_no=${s_three_yes_no:y}
    if [ "${s_three_yes_no}" == "y" ]; then
        aws s3api create-bucket --bucket vault-key-backup-$TF_VAR_app_env --acl private --region $TF_VAR_aws_region --create-bucket-configuration LocationConstraint=$TF_VAR_aws_region
    fi
    log "Create S3 Bucket vault-key-backup in private mode."
fi

log "#############Enable AWS Bucket Encryption############"
log "Check if s3 bucket vault-key-backup-${TF_VAR_app_env} is encrypted."
bucket_check_encrypt=$(aws s3api get-bucket-encryption --region $TF_VAR_aws_region --bucket vault-key-backup-$TF_VAR_app_env | grep AES256)
if [ "$IS_DEBUG" == "true" ]; then
    log "s3 bucket is encrypted check ${bucket_check_encrypt}"
fi
if [ -z "${bucket_check_encrypt}" ]; then
    log "Enable Bucket Encryption"
    aws s3api put-bucket-encryption --bucket vault-key-backup-$TF_VAR_app_env --region $TF_VAR_aws_region --server-side-encryption-configuration '{"Rules": [{"ApplyServerSideEncryptionByDefault": {"SSEAlgorithm": "AES256"}}]}' > /dev/null
fi
