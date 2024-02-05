#!/usr/bin/env bash

version=$1
go_version=$2

git_switch = "!f() { git switch $1 2>/dev/null || git switch -c $1; }; f"

# git_switch gen-v${version}

rm go.mod go.sum
./generate -v ${version} -g ${go_version}

go mod download
go build
