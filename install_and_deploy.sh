#!/bin/bash
set -e
go test ./...
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build .
cp hipchat-openshift docker/http-server
oc start-build openshift-bot --from-dir=docker/
rm docker/http-server
