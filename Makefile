CHAINNET ?= mainnet
BINARY ?= sifnoded
GOPATH ?= $(shell go env GOPATH)
GOBIN ?= $(GOPATH)/bin
NOW = $(shell date +'%Y-%m-%d_%T')
COMMIT := $(shell git log -1 --format='%H')
VERSION := $(shell git describe --match 'v*' --abbrev=8 --tags | sed 's/-g/-/' | sed 's/-[0-9]*-/-/')
IMAGE_TAG ?= latest
HTTPS_GIT := https://github.com/sifchain/sifnode.git
DOCKER ?= docker
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf

GOFLAGS := ""
GOTAGS := ledger

GO_VERSION := $(shell cat go.mod | grep -E 'go [0-9].[0-9]+' | cut -d ' ' -f 2)

LDFLAGS = -X github.com/cosmos/cosmos-sdk/version.Name=sifchain \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=sifnoded \
		  -X github.com/cosmos/cosmos-sdk/version.ClientName=sifnoded \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(GOTAGS)

BUILD_FLAGS := -ldflags '$(LDFLAGS)' -tags '$(GOTAGS)'

BINARIES = ./cmd/sifnoded ./cmd/sifgen ./cmd/ebrelayer ./cmd/siftest

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
	@GOFLAGS=$(GOFLAGS) go mod verify

lint: lint-pre
	@golangci-lint run

lint-verbose: lint-pre
	@golangci-lint run -v --timeout=5m

install: go.sum
	GOFLAGS=$(GOFLAGS) go install $(BUILD_FLAGS) $(BINARIES)

build-sifd: go.sum
	GOFLAGS=$(GOFLAGS) go build  $(BUILD_FLAGS) ./cmd/sifnoded

clean:
	@rm -rf $(GOBIN)/sif*

coverage:
	@GOFLAGS=$(GOFLAGS) go test -v ./... -coverprofile=coverage.txt -covermode=atomic

tests:
	@GOFLAGS=$(GOFLAGS) go test -v -coverprofile .testCoverage.txt ./...

feature-tests:
	@GOFLAGS=$(GOFLAGS) go test -v ./test/bdd --godog.format=pretty --godog.random -race -coverprofile=.coverage.txt

mocks:
	@echo "Generating mocks"

	# Check if mockery is available in $PATH, install it if not.
	@if ! which mockery > /dev/null; then \
		echo "mockery not found, installing version v2..."; \
		go install github.com/vektra/mockery/v2; \
	fi

	# Check if mockgen is available in $PATH, install it if not.
	@if ! which mockgen > /dev/null; then \
		echo "mockgen not found, installing latest version..."; \
		go install go.uber.org/mock/mockgen@latest; \
	fi

	# Check again if mockery is installed, fail if not found.
	@if ! which mockery > /dev/null; then \
		echo "Error: mockery could not be found or installed"; \
		exit 1; \
	fi

	# Check again if mockgen is installed, fail if not found.
	@if ! which mockgen > /dev/null; then \
		echo "Error: mockgen could not be found or installed"; \
		exit 1; \
	fi

	# Run go generate on all packages.
	@go generate ./...

run:
	GOFLAGS=$(GOFLAGS) go run ./cmd/sifnoded start

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

protoVer = v0.3
protoImageName = tendermintdev/sdk-proto-gen:$(protoVer)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName) sh ./scripts/protocgen.sh
.PHONY: proto-gen

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

GORELEASER_IMAGE := ghcr.io/goreleaser/goreleaser-cross:v$(GO_VERSION)

## release: Build binaries for all platforms and generate checksums
ifdef GITHUB_TOKEN
release:
	docker run \
		--rm \
		-e GITHUB_TOKEN=$(GITHUB_TOKEN) \
		-e LDFLAGS="$(LDFLAGS)" \
		-e GOTAGS="$(GOTAGS)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/sifnoded \
		-w /go/src/sifnoded \
		$(GORELEASER_IMAGE) \
		release \
		--clean
else
release:
	@echo "Error: GITHUB_TOKEN is not defined. Please define it before running 'make release'."
endif

## release-dry-run: Dry-run build process for all platforms and generate checksums
release-dry-run:
	docker run \
		--rm \
		-e LDFLAGS="$(LDFLAGS)" \
		-e GOTAGS="$(GOTAGS)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/sifnoded \
		-w /go/src/sifnoded \
		$(GORELEASER_IMAGE) \
		release \
		--clean \
		--skip-publish

## release-snapshot: Build snapshots for all platforms and generate checksums
release-snapshot:
	docker run \
		--rm \
		-e LDFLAGS="$(LDFLAGS)" \
		-e GOTAGS="$(GOTAGS)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/sifnoded \
		-w /go/src/sifnoded \
		$(GORELEASER_IMAGE) \
		release \
		--clean \
		--snapshot \
		--skip-validate \
		--skip-publish

.PHONY: release release-dry-run release-snapshot