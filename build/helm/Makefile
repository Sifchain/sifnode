help: ## Show this help message
	@echo 'usage: make [target] ...'
	@echo
	@echo 'targets:'
	@egrep '^(.+)\:\ ##\ (.+)' ${MAKEFILE_LIST} | column -t -c 2 -s ':#'

NODE?=sifnode
CHAINNET?=testnet

helm: ## Install Helm 3 dependency
	@echo Installing Helm 3 binary
	@curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash

repos: ## Add Helm repositories for dependencies
	@echo Installing Helm repos
	@helm repo add stable https://kubernetes-charts.storage.googleapis.com
	@helm repo add kubernetes-dashboard https://kubernetes.github.io/dashboard

deploy-sifnode: ## Sifnode deployment to Kubernetes
	@echo Sifnode deployment to Kubernetes
	@helm upgrade $(NODE) ./sifnode --set sifnode.env.chainNet=$(CHAINNET) --install -n $(NODE) --create-namespace --debug; \

destroy-sifnode: ## Destroy the Sifnode deployment
	@echo Destroy the Sifnode deployment
	@helm delete $(NODE) -n $(NODE)


.PHONY: deploy-sifnode destroy-sifnode
