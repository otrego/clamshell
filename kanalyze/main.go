// Binary kanalyze anaylzes games, looking for problems.
//
// To process a single sgf:
//
// export KATAGO_MODEL=some_model.bin.gz
// export KATAGO_GTP_CONFIG=some_gtp_config.cfg
// go run main.go foo.sgf
//
// For more about analysis engine, see
// https://github.com/lightvector/KataGo/blob/master/docs/Analysis_Engine.md
package main

import (
	"flag"
	"os"
	"strings"

	"github.com/golang/glog"
)

var output = flag.String("output_dir", "", "Directory for returning the processed SGFs. By default, uses current directory")

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	model := os.Getenv("KATAGO_MODEL")
	gtpConfig := os.Getenv("KATAGO_GTP_CONFIG")

	// TODO(kashomon): Fail if these are empty

	an := &analyzer{
		model:     model,
		gtpConfig: gtpConfig,
	}

	files := getSGFs(flag.Args())

	if err := an.process(files); err != nil {
		glog.Exit(err)
	}
}

func getSGFs(args []string) []string {
	var out []string
	for _, s := range args {
		if strings.HasSuffix(s, ".sgf") {
			out = append(out, s)
		}
		// TODO(kashomon): For directories, read all SGFs in the relevant directory.
	}
	return out
}

type analyzer struct {
	model     string
	gtpConfig string
}

func (z *analyzer) process(files []string) error {
	glog.Infof("using model %q", z.model)
	glog.Infof("using gtp config %q", z.gtpConfig)
	glog.Infof("using files %v\n", files)
	for _, f := range files {
		glog.Infof("analyzing %v\n", f)
		// TODO(kashomon): Actually analyze
	}
	return nil
}
