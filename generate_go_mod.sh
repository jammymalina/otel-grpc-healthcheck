#!/usr/bin/env bash

set -euxo pipefail

version=$1
go_version=$2
component_version=$3

go clean -modcache

rm -f go.mod go.sum
go run scripts/generate_mod_file.go -otelversion="${version}" -goversion="${go_version}" -componentversion="${component_version}

go get ./...
go mod tidy
go vet ./...
go build
