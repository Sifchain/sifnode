CHAINNET?=testnet # Options; testnet, mainnet
BINARY?=sifnoded
GOBIN?=${GOPATH}/bin
NOW=$(shell date +'%Y-%m-%d_%T')
COMMIT:=$(shell git log -1 --format='%H')
VERSION:=$(shell cat version)
TIMESTAMP:=$(shell date +%s)

ifeq (mainnet,${CHAINNET})
	BUILD_TAGS=mainnet
else
	BUILD_TAGS=testnet
endif

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(BUILD_TAGS))

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=sifchain \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=sifnoded \
		  -X github.com/cosmos/cosmos-sdk/version.ClientName=sifnodecli \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"

BUILD_FLAGS := -ldflags '$(ldflags)' -tags $(BUILD_TAGS) -a

BINARIES=./cmd/sifnodecli ./cmd/sifnoded ./cmd/sifgen ./cmd/sifcrg

all: lint install

build-config:
	echo $(CHAINNET)
	echo $(BUILD_TAGS)
	echo $(BUILD_FLAGS)

lint-pre:
	@test -z $(gofmt -l .)
	@go mod verify

lint: lint-pre
	@golangci-lint run

lint-verbose: lint-pre
	@golangci-lint run -v --timeout=5m

install: go.sum
	go install ${BUILD_FLAGS} ${BINARIES}

start:
	sifnodecli rest-server & sifnoded start

clean-config:
	@echo "Are you sure you wish to clear-config? This will destroy your private keys [y/N] " && read ans && [ $${ans:-N} = y ]
	@rm -rf ~/.sifnode*

backup-config:
	mkdir -p .backups
	tar -cvf .backups/sifnode-config-${TIMESTAMP}.tar -C ~/ .sifnoded .sifnodecli

clean:
	@rm -rf ${GOBIN}/sif*

tests:
	@go test -v -coverprofile .testCoverage.txt ./...

feature-tests:
	@go test -v ./test/bdd --godog.format=pretty --godog.random -race -coverprofile=.coverage.txt

build-image:
	docker build -t sifchain/$(BINARY):$(CHAINNET) -f ./cmd/$(BINARY)/Dockerfile .

run-image: build-image
	docker run sifchain/$(BINARY):$(CHAINNET)

sh-image: build-image
	docker run -it sifchain/$(BINARY):$(CHAINNET) sh
