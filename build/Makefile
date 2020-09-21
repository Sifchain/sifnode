include ./build/local/Makefile

NAMESPACE?=dev # Options dev,test,prod
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
