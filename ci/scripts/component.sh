#!/bin/bash -eux

pushd dp-feedback-api
  make test-component
popd
