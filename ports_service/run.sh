#!/bin/bash

set -eou pipefail

PORTS_GRPC_PORT=4040 go run cmd/ports/main.go