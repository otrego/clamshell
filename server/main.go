package main

import (
	"flag"

	glog "github.com/golang/glog"
	"github.com/otrego/clamshell/server/grpc"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

func main() {
	flag.Set("alsologtostderr", "true")
	flag.Parse()

	glog.Info("Starting Clamshell")

	grpc.Run(&grpc.Options{
		Port: *port,
	})
}
