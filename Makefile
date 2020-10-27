CHAINNET?=localnet # Options; localnet, testnet, chaosnet ,mainnet
BINARY?=sifnoded
GOBIN?=${GOPATH}/bin
NOW=$(shell date +'%Y-%m-%d_%T')
COMMIT:=$(shell git log -1 --format='%H')
VERSION:=$(shell cat version)

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=sifchain \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=sifnoded \
		  -X github.com/cosmos/cosmos-sdk/version.ClientName=sifnodecli \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \

BUILD_FLAGS := -ldflags '$(ldflags)' -tags ${CHAINNET} -a

BINARIES=./cmd/sifnodecli ./cmd/sifnoded ./cmd/sifgen ./cmd/sifcrg ./cmd/ebrelayer

all: lint install

lint-pre:
	@test -z $(gofmt -l .)
	@go mod verify

lint: lint-pre
	@golangci-lint run

lint-verbose: lint-pre
	@golangci-lint run -v --timeout=5m

install: go.sum
	go install ${BUILD_FLAGS} ${BINARIES}

clean-config:
	@rm -rf ~/.sifnode*

clean: clean-config
	@rm -rf ${GOBIN}/sif*

tests:
	@go test -v -coverprofile .testCoverage.txt ./...

feature-tests:
	@go test -v ./test/bdd --godog.format=pretty --godog.random -race -coverprofile=.coverage.txt

run:
	go run ./cmd/sifd start


build-image:
	docker build -t sifchain/$(BINARY):$(CHAINNET) -f ./cmd/$(BINARY)/Dockerfile .

run-image: build-image
	docker run sifchain/$(BINARY):$(CHAINNET)

sh-image: build-image
	docker run -it sifchain/$(BINARY):$(CHAINNET) sh
