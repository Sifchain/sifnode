ROOT_DIR = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

CHAIN_ID := "sifchain"
AWS_REGION := "us-west-1"
AWS_CIDR := "10.0.0.0/16"
AWS_CLUSTER_VERSION := 1.17
AWS_CLUSTER_NAME := "sifnode"
AWS_CLUSTER_AMI_TYPE := "AL2_x86_64"
AWS_CLUSTER_DESIRED_CAPACITY := 1
AWS_CLUSTER_MAX_CAPACITY := 1
AWS_CLUSTER_MIN_CAPACITY := 1
AWS_CLUSTER_INSTANCE_TYPE := "t2.small"
AWS_CLUSTER_DISK_SIZE := 100

help: ## Show this help message
	@echo 'usage: make [target] ...'
	@echo
	@echo 'targets:'
	@egrep '^(.+)\:\ ##\ (.+)' ${MAKEFILE_LIST} | column -t -c 2 -s ':#'

gotpl:
	@go get -u github.com/belitre/gotpl

generate-aws-cluster-config: gotpl
	@mkdir -p $(ROOT_DIR)/.environments/providers/aws/eks
	@gotpl $(ROOT_DIR)/template/aws/*.tf.tpl --set chain_id=$(CHAIN_ID) \
		--set aws.region=$(AWS_REGION) \
		--set aws.cidr=$(AWS_CIDR) \
		--set aws.cluster.version=$(AWS_CLUSTER_VERSION) \
		--set aws.cluster.name=$(AWS_CLUSTER_NAME) \
		--set aws.cluster.ami_type=$(AWS_CLUSTER_AMI_TYPE) \
		--set aws.cluster.desired_capacity=$(AWS_CLUSTER_DESIRED_CAPACITY) \
		--set aws.cluster.max_capacity=$(AWS_CLUSTER_MAX_CAPACITY) \
		--set aws.cluster.min_capacity=$(AWS_CLUSTER_MIN_CAPACITY) \
		--set aws.cluster.instance_type=$(AWS_CLUSTER_INSTANCE_TYPE) \
		--set aws.cluster.disk_size=$(AWS_CLUSTER_DISK_SIZE) > $(ROOT_DIR)/.environments/providers/aws/eks/$(CHAIN_ID).tf

build-aws-cluster: generate-aws-cluster-config ## AWS deployment
	cd $(ROOT_DIR)/.environments/providers/aws/eks && terraform init && terraform apply

destroy-aws-cluster: ## Destroy the AWS infrastructure
	cd $(ROOT_DIR)/.environments/providers/aws/eks && terraform destroy

.PHONY: aws aws-destroy
