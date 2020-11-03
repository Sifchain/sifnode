all: run

install:
	pip3 install openapi-spec-validator

run-swagger-ui:
	docker build -f swagger-ui.Dockerfile -t swagger-ui .
	docker run -p 8080:8080 swagger:latest

run-prim:
	docker build -f prism.Dockerfile -t prism .
	docker run -p 1317:1317 prism:latest

test:
	openapi-spec-validator sifchain-openapi.yml

.PHONY: test