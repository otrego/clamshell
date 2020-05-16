package main

import (
	"flag"
	"fmt"
	"net"

	log "github.com/golang/glog"
	pb "github.com/otrego/clamshell/server/api"
	"github.com/otrego/clamshell/server/echo"
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

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterEchoServiceServer(grpcServer, &echo.EchoServer{})
	grpcServer.Serve(lis)
}
