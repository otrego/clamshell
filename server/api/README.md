# API

Protocol buffer installation:

1. `go get google.golang.org/grpc`
2  `go install google.golang.org/protobuf/cmd/protoc-gen-go`
3. Install protobuf via brew: `brew install protobuf`

Generation:

1. To Generate run from root directory

  `protoc -I server/api/ server/api/api.proto --go_out=plugins=grpc:api`

