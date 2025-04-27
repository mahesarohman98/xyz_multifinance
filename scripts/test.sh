#!/bin/bash
set -e

env_file=${1:-.env.test}

ENV_VARS=$(cat "./.env" "./$env_file" | grep -Ev '^#' | xargs)

echo "Running integration tests with env: $env_file"
cd ./src/
env $ENV_VARS go test -count=1 -p=8 -parallel=8 -race ./...

