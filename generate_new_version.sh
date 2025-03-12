#!/usr/bin/env bash

set -euxo pipefail

version=$1
component_version=$2

./generate_go_mod.sh "${version}" "1.24" "${component_version}"

mv go.mod go.mod.new
mv go.sum go.sum.new
mv config.go config.go.new
mv extension.go extension.go.new
mv factory.go factory.go.new
mv README.md README.md.new

git switch main

mv go.mod.new go.mod
mv go.sum.new go.sum
mv config.go.new config.go
mv extension.go.new extension.go
mv factory.go.new factory.go
mv README.md.new README.md

git add go.mod go.sum config.go extension.go factory.go README.md
git commit --allow-empty -m "Updated otel to version ${version}"
git tag -a "v${version}" -m "v${version}"
git push origin HEAD
git push origin "v${version}"
