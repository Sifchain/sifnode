CHAINNET?=betanet
BINARY?=sifnoded
GOPATH?=$(shell go env GOPATH)
GOBIN?=$(GOPATH)/bin
NOW=$(shell date +'%Y-%m-%d_%T')
COMMIT:=$(shell git log -1 --format='%H')
VERSION:=$(shell cat version)
IMAGE_TAG?=latest
HTTPS_GIT := https://github.com/sifchain/sifnode.git
DOCKER ?= docker
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf

GOFLAGS:=""
GOTAGS:=

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=sifchain \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=sifnoded \
		  -X github.com/cosmos/cosmos-sdk/version.ClientName=sifnoded \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

# We use one smart contract file as a signal that the abigen files have been created
smart_contract_file=cmd/ebrelayer/contract/generated/artifacts/contracts/BridgeRegistry.sol/BridgeRegistry.go

BUILD_TAGS ?= ${IMAGE_TAG}
BUILD_FLAGS := -ldflags '$(ldflags)' -tags "$(GOTAGS) ${BUILD_TAGS}"

BINARIES=./cmd/sifnoded ./cmd/sifgen ./cmd/ebrelayer

 # You can regenerate proto_files with
#	find . -name *.proto | sort | grep -v node_mo | grep -v test/integration | xargs echo
# if the list of .proto files changes
#
# go_proto_files is simpler:
#	find . -name *.pb.go | xargs echo

proto_files=./proto/sifnode/admin/v1/query.proto ./proto/sifnode/admin/v1/tx.proto ./proto/sifnode/admin/v1/types.proto ./proto/sifnode/clp/v1/genesis.proto ./proto/sifnode/clp/v1/params.proto ./proto/sifnode/clp/v1/pool.proto ./proto/sifnode/clp/v1/querier.proto ./proto/sifnode/clp/v1/tx.proto ./proto/sifnode/clp/v1/types.proto ./proto/sifnode/dispensation/v1/query.proto ./proto/sifnode/dispensation/v1/tx.proto ./proto/sifnode/dispensation/v1/types.proto ./proto/sifnode/ethbridge/v1/query.proto ./proto/sifnode/ethbridge/v1/tx.proto ./proto/sifnode/ethbridge/v1/types.proto ./proto/sifnode/margin/v1/query.proto ./proto/sifnode/margin/v1/tx.proto ./proto/sifnode/margin/v1/types.proto ./proto/sifnode/oracle/v1/network_descriptor.proto ./proto/sifnode/oracle/v1/types.proto ./proto/sifnode/tokenregistry/v1/query.proto ./proto/sifnode/tokenregistry/v1/tx.proto ./proto/sifnode/tokenregistry/v1/types.proto ./third_party/proto/cosmos/base/coin.proto ./third_party/proto/cosmos/base/query/v1beta1/pagination.proto ./third_party/proto/gogoproto/gogo.proto ./third_party/proto/google/api/annotations.proto ./third_party/proto/google/api/http.proto ./third_party/proto/google/api/httpbody.proto

go_proto_files=./x/ethbridge/types/tx.pb.go ./x/ethbridge/types/types.pb.go ./x/ethbridge/types/query.pb.go ./x/oracle/types/network_descriptor.pb.go ./x/oracle/types/types.pb.go ./x/admin/types/tx.pb.go ./x/admin/types/types.pb.go ./x/admin/types/query.pb.go ./x/tokenregistry/types/tx.pb.go ./x/tokenregistry/types/types.pb.go ./x/tokenregistry/types/query.pb.go ./x/margin/types/tx.pb.go ./x/margin/types/types.pb.go ./x/margin/types/query.pb.go ./x/clp/types/pool.pb.go ./x/clp/types/tx.pb.go ./x/clp/types/querier.pb.go ./x/clp/types/genesis.pb.go ./x/clp/types/params.pb.go ./x/dispensation/types/tx.pb.go ./x/dispensation/types/types.pb.go ./x/dispensation/types/query.pb.go


.PHONY: all
all: lint install

build-config:
	echo $(CHAINNET)
	echo $(BUILD_FLAGS)

init:
	./scripts/init.sh

start:
	sifnoded start

# Note that ebrelayer depends on go files from the smart contracts, so the smart contracts
# must be built first
lint-pre: ${smart_contract_file}
	# test -z "$(shell gofmt -l .)"
	GOFLAGS=${GOFLAGS} go mod verify

lint: lint-pre
	golangci-lint run

lint-verbose: lint-pre
	golangci-lint run -v --timeout=5m

.PHONY: install
install: ${BINARIES}

${BINARIES} &: go.mod go.sum ${smart_contract_file} $(go_proto_files)
	GOFLAGS=${GOFLAGS} go install ${BUILD_FLAGS} ${BINARIES}
	# You can't depend on go updating the timestamps - go install may decide it doesn't need to do any work
	touch ${BINARIES}

.PHONY: install-smart-contracts
install-smart-contracts: ${smart_contract_file}
	make -C smart-contracts install

clean-config:
	@rm -rf ~/.sifnode*

clean-ebrelayer:
	@rm -rf ${GOBIN}/ebrelayer

clean-peggy:
	make -C smart-contracts clean

clean: clean-config clean-peggy clean-ebrelayer
	@rm -rf ${GOBIN}/sif*
	git clean -fdx cmd/ebrelayer/contract/generated

coverage:
	GOFLAGS=${GOFLAGS} go test -v ./... -coverprofile=coverage.txt -covermode=atomic

.PHONY: tests test-peggy test-bin feature-tests
test-peggy:
	$(MAKE) -C smart-contracts tests

test-bin: ${BINARIES}
	GOFLAGS=${GOFLAGS} go test -v -coverprofile .testCoverage.txt ./...

.PHONY: tests test
test tests: test-peggy test-bin

feature-tests:
	GOFLAGS=${GOFLAGS} go test -v ./test/bdd --godog.format=pretty --godog.random -race -coverprofile=.coverage.txt

run:
	GOFLAGS=$(GOFLAGS) go run ./cmd/sifnoded start

build-image: install-smart-contracts
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

protoVer=v0.3
protoImageName=tendermintdev/sdk-proto-gen:$(protoVer)

proto-all: proto-format proto-lint $(go_proto_files)

.PHONY: proto-gen
proto-gen: $(go_proto_files)

$(go_proto_files) &: $(proto_files)
	@echo "Generating Protobuf files"
	$(DOCKER) run -e SIFUSER=$(shell id -u):$(shell id -g) --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/sdk-proto-gen:v0.3 sh -x ./scripts/protocgen.sh
	touch $@

proto-format:
	@echo "Formatting Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace \
	--workdir /workspace $(protoImageName) \
	find ./ -not -path "./third_party/*" -name *.proto -exec clang-format -i {} \;
.PHONY: proto-format

# This generates the SDK's custom wrapper for google.protobuf.Any. It should only be run manually when needed
proto-gen-any:
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName) sh ./scripts/protocgen-any.sh
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

${smart_contract_file}:
	$(MAKE) -C smart-contracts
