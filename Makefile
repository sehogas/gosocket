createnetwork:
	docker network create -d bridge go-network

run:
	go run main.go

build:
	docker build -t gosocket:1.0 .

bin:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o gosocket .

install:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags="-s -w"

tag:
	docker tag gosocket:1.0 gosocket:latest

start:
	docker run --rm -d --name gosocket -p 8000:8000 gosocket:latest

start_old:
	docker run --rm -d --name gosocket --network go-network -p 8000:8000 gosocket:latest

stop:
	docker stop gosocket

.PHONY: createnetwork run build bin install tag start stop  
