#!/bin/bash -eux

go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

pushd dp-feedback-api
  make lint
popd
