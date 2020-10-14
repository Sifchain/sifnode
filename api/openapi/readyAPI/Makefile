build:
	docker build -t readyapi .

run: build
	docker run -p 9090:9090 -p 8000-8010:8000-8010 readyapi:latest
