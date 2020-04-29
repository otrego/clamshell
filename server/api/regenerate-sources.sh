#!/bin/bash

set -e

cd "$(dirname "$0")"

# Change to top-level directory.
cd ../../

DIR=$PWD

protoc --go_out=plugins=grpc,paths=source_relative:. \
  -I . \
  ./server/api/*.proto
