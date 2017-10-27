#!/bin/bash
set -e
go test ./...
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build .
mv http-server docker/
oc start-build openshift-bot --from-dir=docker/
rm docker/http-server