#!/bin/bash -eux

pushd dp-feedback-api
  make build
  cp build/dp-feedback-api Dockerfile.concourse ../build
popd
