CHAINNET?=betanet
BINARY?=sifnoded
GOBIN?=${GOPATH}/bin
NOW=$(shell date +'%Y-%m-%d_%T')
COMMIT:=$(shell git log -1 --format='%H')
VERSION:=$(shell cat version)
IMAGE_TAG?=latest
HTTPS_GIT := https://github.com/sifchain/sifnode.git
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=sifchain \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=sifnoded \
		  -X github.com/cosmos/cosmos-sdk/version.ClientName=sifnoded \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

BINARIES=./cmd/sifnoded ./cmd/sifgen ./cmd/ebrelayer

all: lint install

build-config:
	echo $(CHAINNET)
	echo $(BUILD_FLAGS)

init:
	./scripts/init.sh

start:
	sifnoded start

lint-pre:
	@test -z $(gofmt -l .)
	@go mod verify

lint: lint-pre
	@golangci-lint run

lint-verbose: lint-pre
	@golangci-lint run -v --timeout=5m

install-all: go.sum
	go install ${BUILD_FLAGS} ${BINARIES}

install: install-all

install-ebrelayer: go.sum
	go install ${BUILD_FLAGS} ./cmd/ebrelayer

install-sifgen: go.sum
	go install ${BUILD_FLAGS} ./cmd/sifgen

install-sifnode: go.sum
	go install ${BUILD_FLAGS} ./cmd/sifgen

install-sifnode-rocks: go.sum
	go install ${BUILD_FLAGS} -tags rocksdb ./cmd/sifnoded

build-sifd: go.sum
	go build  ${BUILD_FLAGS} ./cmd/sifnoded

clean:
	@rm -rf ${GOBIN}/sif*

coverage:
	@go test -v ./... -coverprofile=coverage.txt -covermode=atomic

tests:
	@go test -v -coverprofile .testCoverage.txt ./...

feature-tests:
	@go test -v ./test/bdd --godog.format=pretty --godog.random -race -coverprofile=.coverage.txt

run:
	go run ./cmd/sifnoded start

build-image:
	docker build -t sifchain/$(BINARY):$(IMAGE_TAG) -f ./cmd/$(BINARY)/Dockerfile .

run-image: build-image
	docker run sifchain/$(BINARY):$(IMAGE_TAG)

sh-image: build-image
	docker run -it sifchain/$(BINARY):$(IMAGE_TAG) sh

init-run:
	./scripts/init.sh && ./scripts/run.sh

init-run-noInstall:
	./scripts/init-noInstall.sh && ./scripts/run.sh

rollback:
	./scripts/rollback.sh

###############################################################################
###                                Protobuf                                 ###
###############################################################################

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/sdk-proto-gen sh ./scripts/protocgen.sh
.PHONY: proto-gen

proto-format:
	@echo "Formatting Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace \
	--workdir /workspace tendermintdev/docker-build-proto \
	find ./ -not -path "./third_party/*" -name *.proto -exec clang-format -i {} \;
.PHONY: proto-format

# This generates the SDK's custom wrapper for google.protobuf.Any. It should only be run manually when needed
proto-gen-any:
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/sdk-proto-gen sh ./scripts/protocgen-any.sh
.PHONY: proto-gen-any

proto-swagger-gen:
	@./scripts/protoc-swagger-gen.sh
.PHONY: proto-swagger-gen

proto-lint:
	$(DOCKER_BUF) lint --error-format=json
.PHONY: proto-lint

proto-check-breaking:
	# we should turn this back on after our first release
	# $(DOCKER_BUF) breaking --against $(HTTPS_GIT)#branch=master
.PHONY: proto-check-breaking
