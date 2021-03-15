CHAINNET?=testnet # Options; localnet, testnet, chaosnet ,mainnet
BINARY?=sifnoded
GOBIN?=${GOPATH}/bin
NOW=$(shell date +'%Y-%m-%d_%T')
COMMIT:=$(shell git log -1 --format='%H')
VERSION:=$(shell cat version)

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

BUILD_FLAGS := -ldflags '$(ldflags)' -tags ${BUILD_TAGS} -a

BINARIES=./cmd/sifnodecli ./cmd/sifnoded ./cmd/sifgen ./cmd/sifcrg ./cmd/ebrelayer

all: lint install

build-config:
	echo $(CHAINNET)
	echo $(BUILD_TAGS)
	echo $(BUILD_FLAGS)

init:
	./scripts/init.sh

start:
	sifnodecli rest-server & sifnoded start

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

coverage:
	@go test -v ./... -coverprofile=coverage.txt -covermode=atomic

tests:
	@go test -v -coverprofile .testCoverage.txt ./...

feature-tests:
	@go test -v ./test/bdd --godog.format=pretty --godog.random -race -coverprofile=.coverage.txt

run:
	go run ./cmd/sifnoded start

build-image:
	docker build -t sifchain/$(BINARY):$(IMAGE_TAG) --build-arg chainnet=$(CHAINNET) -f ./cmd/$(BINARY)/Dockerfile .

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

