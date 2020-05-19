// Package grpc contains helpers for initializing gRPC
package grpc

import (
	"fmt"
	"net"

	log "github.com/golang/glog"
	"google.golang.org/grpc"

	pb "github.com/otrego/clamshell/server/api"
	"github.com/otrego/clamshell/server/echo"
)

// Options contains options for running gRPC
type Options struct {
	// Port to listen on
	Port int
}

// Run starts an in-process gRPC gateway and gRPC server.
//
// The server will be shutdown when "ctx" is canceled.
func Run(opts *Options) {
	addr := fmt.Sprintf("localhost:%d", opts.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to dial to addr %s: %v", addr, err)
	}
	grpcServer := initGRPCServer(opts)
	log.Infof("listening on :%d", opts.Port)
	grpcServer.Serve(lis)
}

func initGRPCServer(_ *Options) *grpc.Server {
	var gopts []grpc.ServerOption
	grpcServer := grpc.NewServer(gopts...)
	pb.RegisterEchoServiceServer(grpcServer, &echo.EchoServer{})
	return grpcServer
}
