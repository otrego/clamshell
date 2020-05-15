#!/bin/bash

set -e

cd "$(dirname "$0")"

DIR=$PWD

protoc --go_out=plugins=grpc,paths=source_relative:. \
  -I . \
  -I ./third_party/googleapis \
  -I ./third_party/grpc-gateway \
  ./server/api/*.proto

protoc --grpc-gateway_out=logtostderr=true,paths=source_relative:. \
  -I . \
  -I ./third_party/googleapis \
  -I ./third_party/grpc-gateway \
  ./server/api/*.proto
