include ./build/Makefile

CHAINNET?=localnet # Options; localnet, testnet, chaosnet ,mainnet
GOBIN?=${GOPATH}/bin
NOW=$(shell date +'%Y-%m-%d_%T')
COMMIT:=$(shell git log -1 --format='%H')
VERSION:=$(shell cat version)


ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=SifChain \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=sifd \
		  -X github.com/cosmos/cosmos-sdk/version.ClientName=sifcli \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(TAG)

BUILD_FLAGS := -ldflags '$(ldflags)' -tags ${TAG} -a

BINARIES=./cmd/sifnodecli ./cmd/sifnoded

install: go.sum
	go install ${BUILD_FLAGS} ${BINARIES}


feature-tests:
	@go test -v ./test/bdd --godog.format=pretty --godog.random -race -coverprofile=.coverage.txt

run:
	go run ./cmd/sifd start

