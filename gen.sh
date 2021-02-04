#!/bin/bash

protoc ./ports.proto --go_out=./ports --go-grpc_out=./ports
