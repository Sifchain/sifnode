help: ## Show this help message
	@echo 'usage: make [target] ...'
	@echo
	@echo 'targets:'
	@egrep '^(.+)\:\ ##\ (.+)' ${MAKEFILE_LIST} | column -t -c 2 -s ':#'

helm: ## Install Helm 3 dependency
	@echo Installing Helm 3 binary
	@curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash

repos: ## Add Helm repositories for dependencies
	@echo Installing Helm repos
	@helm repo add stable https://kubernetes-charts.storage.googleapis.com
	@helm repo add kubernetes-dashboard https://kubernetes.github.io/dashboard

deploy-sifnode: ## Sifnode deployment to Kubernetes
	@echo Sifnode deployment to Kubernetes
	@helm upgrade sifnode ./sifnode --install -n sifnode --create-namespace --wait

destroy-sifnode: ## Destroy the Sifnode deployment
	@echo Destroy the Sifnode deployment
	@helm delete sifnode -n sifnode

.PHONY: deploy-sifnode destroy-sifnode
