include ./build/local/Makefile

CHAINNET?=localnet # Options; localnet, testnet, chaosnet ,mainnet
NAMESPACE?=dev # Optional dev,test,prod, node0, node1
PROVIDER?=local# local, aws, gcp, az, do
TAG=$(CHAINNET)

build-image:
	docker build -t sifnode:$(TAG) .

run-image: build-image
	docker run sifnode:$(TAG)

build-cluster:
ifeq '$(PROVIDER)' 'local'
	echo "local setup..."
	Make microk8s-cluster-start
endif
