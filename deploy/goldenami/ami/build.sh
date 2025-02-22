#!/usr/bin/env bash

set -euo pipefail

echo "Building frontend..."
npm run --prefix=../../../frontend build

echo "Building backend..."
pushd ../../../backend/
make build
popd

echo "Building AMI..."
packer build .
