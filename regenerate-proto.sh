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

protoc --swagger_out=logtostderr=true:. \
  -I . \
  -I ./third_party/googleapis \
  -I ./third_party/grpc-gateway \
  ./server/api/*.proto

# Manually go-embed the swagger into a go file.
#
# This is quite hacky -- we should use go-bindata or even go generate. For now,
# it's mega simple. See https://github.com/otrego/clamshell/issues/40
SWAG_GO="./server/api/api.swagger.go"
echo 'package api' > "${SWAG_GO}"
echo '' >> "${SWAG_GO}"
echo '// Swagger contains embedded swagger data.' >> "${SWAG_GO}"
echo 'const Swagger = `' >> "${SWAG_GO}"
cat ./server/api/api.swagger.json | sed 's/`/'"'"'/g' >> "${SWAG_GO}"
echo '`' >> "${SWAG_GO}"
