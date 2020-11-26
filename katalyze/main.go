// Binary katalyze anaylzes games, looking for problems.
//
// To process a single sgf:
//
// go run katalyze/main.go \
// -model=./katalyze/testdata/g170e-b10c128-s1141046784-d204142634.bin.gz \
// -config=./katalyze/testdata/analysis_example.cfg \
// katalyze/testdata/example-game.sgf
//
//
// OR
//
// export KATAGO_MODEL=some_model.bin.gz
// export KATAGO_GTP_CONFIG=some_gtp_config.cfg
// go run main.go foo.sgf
//
// To get higher visibility into logs:
//
// Show katago stderr:
// go run main.go -v=2
//
// Show very detailed logs in katago:
// go run main.go -v=3
//
// For more about analysis engine, see
// https://github.com/lightvector/KataGo/blob/master/docs/Analysis_Engine.md
package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/otrego/clamshell/core/game"
	"github.com/otrego/clamshell/core/katago"
	"github.com/otrego/clamshell/core/katago/kataprob"
	"github.com/otrego/clamshell/core/sgf"
)

var (
	outputDir       = flag.String("output_dir", "", "Directory for returning the processed SGFs. By default, uses current directory")
	modelFlag       = flag.String("model", "", "The model to use for katago. If not set, looks for env var $KATAGO_MODEL. Example: g170-b10c128-s197428736-d67404019.bin.gz")
	configFlag      = flag.String("config", "", "The analysis config file to use. If not set, looks for env var $KATAGO_ANALYSIS_CONFIG. Example: analysis_example.cfg")
	analysisThreads = flag.Int("analysis_threads", 8, "The number of analysis threads")

	startFromMove = flag.Int("start_from_move", 0, "Start all games from the specified move number; this is primarily used for debugging")
	maxMoves      = flag.Int("max_moves_per_game", 0, "Only allow this many moves per game to be analyzed. If specified at zero, analyze the whole game")
)

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()

	model := os.Getenv("KATAGO_MODEL")
	if *modelFlag != "" {
		model = *modelFlag
	}

	if model == "" {
		glog.Exit("--model=<katago-model> must be specified, but was empty")
	}

	analysisConfig := os.Getenv("KATAGO_ANALYSIS_CONFIG")
	if *configFlag != "" {
		analysisConfig = *configFlag
	}
	if analysisConfig == "" {
		glog.Exit("--config=<analysis-config> must be specified, but was empty")
	}

	files := getSGFs(flag.Args())
	if len(files) == 0 {
		glog.Exit("No SGF files specified -- must be specified as args to katalyze")
	}

	an := katago.New(model, analysisConfig, *analysisThreads)

	if err := an.Start(); err != nil {
		glog.Exitf("error booting Katago: %v", err)
	}

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
		glog.Infof("Processing file %q", fi)

		content, err := ioutil.ReadFile(fi)
		if err != nil {
			return err
		}
		g, err := sgf.FromString(string(content)).Parse()
		if err != nil {
			return err
		}
		q, err := katago.AnalysisQueryFromGame(g, &katago.QueryOptions{
			MaxMoves:  maxMoves,
			StartFrom: startFromMove,
		})
		if err != nil {
			return err
		}
		result, err := an.AnalyzeGame(q)
		if err != nil {
			return err
		}
		glog.Infof("Finished processing file %q", fi)
		glog.V(2).Infof("Processing data: %v\n", result)

		err = result.AddToGame(g)
		if err != nil {
			return err
		}
		glog.Infof("Finished adding to game for file %q", fi)

		var positions []game.Treepath
		if paths, err := kataprob.FindBlunders(g); err != nil {
			return err
		} else {
			positions = append(positions, paths...)
		}
		var found []string
		for _, tp := range positions {
			found = append(found, tp.CompactString())
		}
		glog.Infof("Found Positions: %v", found)
	}
	return nil
}
