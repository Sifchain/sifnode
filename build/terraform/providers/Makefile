ROOT_DIR = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

CHAINNET := "localnet"
STATE_BACKEND := "local"
AWS_REGION := "us-west-1"
AWS_STATE_BUCKET_NAME := "state"
AWS_CIDR := "10.0.0.0/16"
AWS_CLUSTER_VERSION := 1.17
AWS_CLUSTER_NAME := "sifnode"
AWS_CLUSTER_AMI_TYPE := "AL2_x86_64"
AWS_CLUSTER_DESIRED_CAPACITY := 1
AWS_CLUSTER_MAX_CAPACITY := 1
AWS_CLUSTER_MIN_CAPACITY := 1
AWS_CLUSTER_INSTANCE_TYPE := "t2.small"
AWS_CLUSTER_DISK_SIZE := 100

gotpl:
	@go get -u github.com/belitre/gotpl

create-new-aws-environment: gotpl
	@mkdir -p .live/sifchain-aws-$(CHAINNET)
	@gotpl $(ROOT_DIR)/template/aws/state/$(STATE_BACKEND).tf.tpl $(ROOT_DIR)/template/aws/cluster.tf.tpl \
		--set chainnet=$(CHAINNET) \
		--set aws.state.bucket_name=$(AWS_STATE_BUCKET_NAME) \
		--set aws.state.bucket_key=sifchain-aws-$(CHAINNET).state \
		--set aws.region=$(AWS_REGION) \
		--set aws.cidr=$(AWS_CIDR) \
		--set aws.cluster.version=$(AWS_CLUSTER_VERSION) \
		--set aws.cluster.name=$(AWS_CLUSTER_NAME) \
		--set aws.cluster.ami_type=$(AWS_CLUSTER_AMI_TYPE) \
		--set aws.cluster.desired_capacity=$(AWS_CLUSTER_DESIRED_CAPACITY) \
		--set aws.cluster.max_capacity=$(AWS_CLUSTER_MAX_CAPACITY) \
		--set aws.cluster.min_capacity=$(AWS_CLUSTER_MIN_CAPACITY) \
		--set aws.cluster.instance_type=$(AWS_CLUSTER_INSTANCE_TYPE) \
		--set aws.cluster.disk_size=$(AWS_CLUSTER_DISK_SIZE) > .live/sifchain-aws-$(CHAINNET)/main.tf

.PHONY: aws aws-destroy
