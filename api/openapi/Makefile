all: run

install:
	pip3 install openapi-spec-validator

build:
	docker build -t swagger .

run: build
	docker run -p 8080:8080 swagger:latest

test:
	openapi-spec-validator --schema 2.0 swagger.yml
	openapi-spec-validator clp_get.yml
	openapi-spec-validator clp_post.yml

.PHONY: test