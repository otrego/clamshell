package main

import (
	"flag"

	glog "github.com/golang/glog"
	"github.com/otrego/clamshell/server/config"
	"github.com/otrego/clamshell/server/grpc"
)

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	glog.Info("Starting Clamshell")

	spec, err := config.FromEnv()
	if err != nil {
		glog.Exit(err)
	}

	grpc.Run(spec)
}
