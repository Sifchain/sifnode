CHAINNET?=localnet # Options; localnet, testnet, chaosnet ,mainnet

all: run

install:
	pip3 install openapi-spec-validator

run-swagger-ui: build-image-swagger-iu
	docker run -p 8080:8080 sifchain/swagger-ui:$(CHAINNET)

run-prism: build-image-prism
	docker run -p 1317:1317 sifchain/prism:$(CHAINNET)

build-image-swagger-iu:
	docker build -f swagger-ui.Dockerfile -t sifchain/swagger-ui:$(CHAINNET) .

build-image-prism:
	docker build -f prism.Dockerfile -t sifchain/prism:$(CHAINNET) .

test:
	openapi-spec-validator sifchain-openapi.yml

.PHONY: test