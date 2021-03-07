#!/bin/bash

protoc ./ports.proto --go_out=./ports_service/ports --go-grpc_out=./ports_service/ports \
  --go_out=./client_service/ports --go-grpc_out=./client_service/ports
