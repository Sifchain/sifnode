help: ## Show this help message
	@echo 'usage: make [target] ...'
	@echo
	@echo 'targets:'
	@egrep '^(.+)\:\ ##\ (.+)' ${MAKEFILE_LIST} | column -t -c 2 -s ':#'

NODE=$(nodename)

helm: ## Install Helm 3 dependency
	@echo Installing Helm 3 binary
	@curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash

repos: ## Add Helm repositories for dependencies
	@echo Installing Helm repos
	@helm repo add stable https://kubernetes-charts.storage.googleapis.com
	@helm repo add kubernetes-dashboard https://kubernetes.github.io/dashboard

deploy-sifnode: ## Sifnode deployment to Kubernetes
	@echo Sifnode deployment to Kubernetes
	@if [[ $NODE == '' ]]; then \
		echo "The 'nodename' argument is not defined"; \
		exit 1; \
	fi
	@helm upgrade $(NODE) ./sifnode --install -n $(NODE) --create-namespace --wait; \

destroy-sifnode: ## Destroy the Sifnode deployment
	@echo Destroy the Sifnode deployment
	@if [[ $NODE == '' ]]; then \
		echo "The 'nodename' argument is not defined"; \
		exit 1; \
	fi
	@helm delete $(NODE) -n $(NODE); \


.PHONY: deploy-sifnode destroy-sifnode
