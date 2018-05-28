#!/usr/bin/env bash

set -o errexit
set -o nounset

timeout=60s

for d in $(go list ./...); do
    go test -v -timeout $timeout -tags $1 -race $d
done
exit
