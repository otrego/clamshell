package main

import (
	"flag"

	glog "github.com/golang/glog"
	"github.com/otrego/clamshell/server/config"
	"github.com/otrego/clamshell/server/serve"
)

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	glog.Info("Starting Clamshell")

	spec, err := config.FromEnv()
	if err != nil {
		glog.Exit(err)
	}

	serve.Run(spec)
}
