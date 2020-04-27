# API

Protocol buffer installation:

1. `go get google.golang.org/grpc`
2. `go get github.com/golang/protobuf`
2. `go install github.com/golang/protobuf/protoc-gen-go`
3. Install protobuf compiler via brew:
    1. Mac: `brew install protobuf`


Note:  Don't install: `google.golang.org/protobuf/cmd/protoc-gen-go`. This
doesn't support gRPC yet. See: https://stackoverflow.com/questions/60578892/protoc-gen-go-grpc-program-not-found-or-is-not-executable

Generation:

1. To Generate go-code run:

  `./regenerate-sources.sh`
