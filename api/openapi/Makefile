CHAINNET?=localnet # Options; localnet, testnet, chaosnet ,mainnet

all: run

install:
	pip3 install openapi-spec-validator

run-swagger-ui: build-image-swagger-iu
	docker run -p 8080:8080 sifchain/swagger-ui:$(CHAINNET)

run-prim: build-image-prim
	docker run -p 1317:1317 sifchain/prism:$(CHAINNET)

build-image-swagger-iu:
	docker build -f swagger-ui.Dockerfile -t sifchain/swagger-ui:$(CHAINNET) .

build-image-prim:
	docker build -f prism.Dockerfile -t sifchain/prism:$(CHAINNET) .

test:
	openapi-spec-validator sifchain-openapi.yml

.PHONY: test