// Binary katalyze anaylzes games, looking for problems.
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
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/otrego/clamshell/core/katago"
	"github.com/otrego/clamshell/core/sgf"
)

var (
	outputDir  = flag.String("output_dir", "", "Directory for returning the processed SGFs. By default, uses current directory")
	modelFlag  = flag.String("model", "", "The model to use for katago. If not set, looks for env var $KATAGO_MODEL. Example: g170-b10c128-s197428736-d67404019.bin.gz")
	configFlag = flag.String("config", "", "The analysis config file to use. If not set, looks for env var $KATAGO_ANALYSIS_CONFIG. Example: analysis_example.cfg")
)

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	model := os.Getenv("KATAGO_MODEL")
	if *modelFlag != "" {
		model = *modelFlag
	}

	analysisConfig := os.Getenv("KATAGO_ANALYSIS_CONFIG")
	if *configFlag != "" {
		analysisConfig = *configFlag
	}

	an := katago.GetAnalyzer(model, analysisConfig, 0)
	err := an.Start()

	if err != nil {
		glog.Exitf("error booting Katago: %v", err)
	}

	files := getSGFs(flag.Args())

	if err := process(files, an); err != nil {
		glog.Exitf("error processing files: %v", err)
	}

	an.Stop()
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

func process(files []string, an *katago.Analyzer) error {
	glog.Infof("using files %v\n", files)
	for _, fi := range files {
		content, err := ioutil.ReadFile(fi)
		if err != nil {
			return err
		}
		game, err := sgf.FromString(string(content)).Parse()
		if err != nil {
			return err
		}
		q, err := katago.MainBranchSurvey(game)
		if err != nil {
			return err
		}
		jsonStr, err := q.ToJSON()
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", string(jsonStr))

		result, err := an.AnalyzeGame(q)

		if err != nil {
			return err
		}

		fmt.Printf("%v\n", result)
	}
	return nil
}
