#!/usr/bin/env bash

set -euxo pipefail

version=$1
go_version=$2

go clean -modcache

rm -f go.mod go.sum
./generate.js -o ${version} -g ${go_version}

go get ./...
go mod tidy
go vet ./...
go build
