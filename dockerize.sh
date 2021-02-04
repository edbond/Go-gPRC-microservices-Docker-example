#!/bin/bash

docker build -f ports_service/Dockerfile -t ports_service .
docker build -f client_service/Dockerfile -t client_service .