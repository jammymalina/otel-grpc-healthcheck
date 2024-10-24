!/usr/bin/env bash

set -euxo pipefail

v=$1

./generate_go_mod.sh ${v} "1.23"

mv go.mod go.mod.new
mv go.sum go.sum.new
mv config.go config.go.new
mv extension.go extension.go.new
mv factory.go factory.go.new

git switch main

mv go.mod.new go.mod
mv go.sum.new go.sum
mv config.go.new config.go
mv extension.go.new extension.go
mv factory.go.new factory.go

git add go.mod go.sum config.go extension.go factory.go
git commit --allow-empty -m "Updated otel to version ${v}"
git tag -a "v${v}" -m "v${v}"
git push origin HEAD
git push origin "v${v}"
