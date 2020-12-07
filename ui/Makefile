CHAINNET?=localnet # Options; localnet, testnet, chaosnet ,mainnet

run-ui: build-ui
	docker run -p 8080:80 sifchain/ui:$(CHAINNET)

build-ui:
	docker build --build-arg SIFNODE_API=https://api-monkeybars.sifchain.finance/ -f Dockerfile -t sifchain/ui:$(CHAINNET) .

.PHONY: run-ui build-ui
