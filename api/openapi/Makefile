CHAINNET?=localnet # Options; localnet, testnet, chaosnet ,mainnet

all: run

install:
	pip3 install openapi-spec-validator

run-swagger-ui: build-swagger-iu
	docker run -p 8080:8080 swagger:latest

run-prim: build-prim
	docker run -p 1317:1317 prism:latest

build-image-swagger-iu:
	docker build -f swagger-ui.Dockerfile -t sifchain/swagger-ui:$(CHAINNET) .

build-image-prim:
	docker build -f prism.Dockerfile -t sifchain/prism:$(CHAINNET) .

test:
	openapi-spec-validator sifchain-openapi.yml

.PHONY: test