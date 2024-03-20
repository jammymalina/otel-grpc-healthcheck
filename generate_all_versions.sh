!/usr/bin/env bash

set -euxo pipefail

versions=("0.96.0")

for v in ${versions[@]}; do
  git switch generate-otel-version
  ./generate_go_mod.sh ${v} "1.20"
  mv go.mod go.mod.new
  mv go.sum go.sum.new
  git switch main
  mv go.mod.new go.mod
  mv go.sum.new go.sum
  git add go.mod go.sum
  git commit --allow-empty -m "Updated otel to version ${v}"
  git tag -a "v${v}" -m "v${v}"
  git push origin HEAD
  git push origin "v${v}"
done

