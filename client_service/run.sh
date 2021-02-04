#!/bin/bash

export PORTS_ADDRESS=localhost:4040 
export PORTS_JSON=../ports/ports.json 
export HTTP_PORT=8080 

go run cmd/client/main.go
# docker run --rm -it -p 8080:8080 client_service