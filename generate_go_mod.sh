#!/usr/bin/env bash

set -euxo pipefail

version=$1
component_version=$2

rm -f go.mod go.sum

pushd scripts

go clean -modcache
otel_go_mod_file="otel.go.mod"
curl -o "${otel_go_mod_file}" "https://raw.githubusercontent.com/open-telemetry/opentelemetry-collector-contrib/v${version}/go.mod"
go run generate_mod_file.go -otelversion="${version}" -otelgomodpath="${otel_go_mod_file}" -componentversion="${component_version}"

popd

go clean -modcache

go get ./...
go mod tidy
go vet ./...
go build
