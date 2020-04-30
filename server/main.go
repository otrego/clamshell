package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	log "github.com/golang/glog"
	pb "github.com/otrego/clamshell/server/api"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 10000, "The server port")
)

func main() {
	flag.Set("alsologtostderr", "true")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Infof("listening on :%d", *port)

	// TODO(kashomon): Add TLS
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterFooServiceServer(grpcServer, &fooServer{})
	grpcServer.Serve(lis)
}

type fooServer struct {
}

// GetFoo gets a foo
func (s *fooServer) GetFoo(ctx context.Context, req *pb.FooRequest) (*pb.Foo, error) {
	return &pb.Foo{Content: "zed"}, nil
}

// ListFeatures lists all features contained within the given bounding Rectangle.
func (s *fooServer) ListFoo(ctx context.Context, _ *pb.EmptyRequest) (*pb.FooCollection, error) {
	return &pb.FooCollection{
		Foos: []*pb.Foo{
			{Content: "zed"},
			{Content: "zorp"},
		},
	}, nil
}
