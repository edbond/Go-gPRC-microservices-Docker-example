#!/bin/bash

set -e

for folder in client_service ports_service; do
    pushd $folder
    go test ./...
    popd
done