package main

import (
	"flag"

	_ "github.com/golang/glog"
	"github.com/otrego/clamshell/server/grpc"
)

var (
	port = flag.Int("port", 10000, "The server port")
)

func main() {
	flag.Set("alsologtostderr", "true")
	flag.Parse()

	grpc.Run(&grpc.Options{
		Port: *port,
	})
}
