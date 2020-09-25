help: ## Show this help message
	@echo 'usage: make [target] ...'
	@echo
	@echo 'targets:'
	@egrep '^(.+)\:\ ##\ (.+)' ${MAKEFILE_LIST} | column -t -c 2 -s ':#'

aws: ## AWS deployment
	cd aws && terraform init && terraform apply

aws-destroy: ## Destroy the AWS infrastructure
	cd aws && terraform destroy

.PHONY: aws aws-destroy
