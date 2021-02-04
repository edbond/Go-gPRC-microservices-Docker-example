#!/bin/bash

set -e

for folder in client_service ports_service ports; do
    pushd $folder
    go test ./...
    popd
done