include ./build/terraform/providers/Makefile

CHAINNET?=localnet # Options; localnet, testnet, chaosnet ,mainnet
NAMESPACE?=dev # Optional dev,test,prod, node0, node1
PROVIDER?=local# local, aws, gcp, az, do
BINARY?=sifnoded

build-image:
	docker build -t sifchain/$(BINARY):$(CHAINNET) -f cmd/$(BINARY)/Dockerfile .

run-image: build-image
	docker run sifchain/$(BINARY):$(CHAINNET)

build-cluster:
ifeq '$(PROVIDER)' 'local'
	echo "local setup..."
	make microk8s-cluster-start
endif

create-new-environment:
ifeq '$(PROVIDER)' 'aws'
	make create-new-aws-environment
endif
