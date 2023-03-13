#!/bin/bash

set -eou pipefail

for folder in client_service ports_service; do
    pushd $folder
    go test ./...
    popd
done