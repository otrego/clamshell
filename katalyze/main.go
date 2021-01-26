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
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/golang/glog"
	"github.com/otrego/clamshell/core/katago"
	"github.com/otrego/clamshell/core/storage"
)

var (
	storageProvider = flag.String("storage_provider", "", "Storage provider where analyzed games will be stored. Optional {localdisk|gcp}. Default is localdisk")
	outputDir       = flag.String("output_dir", "", "Directory for returning the processed SGFs. By default, uses current directory")
	bucket          = flag.String("bucket", "", "The bucket for CloudStorage of analyzed games. This param is optional.")
	bucketPrefix    = flag.String("bucket_prefix", "", "Dictates the path within the bucket where games wil be stored. Required if the bucket parameter is supplied.")
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

	files, err := filterSGFs(flag.Args())
	if err != nil {
		glog.Exit(err)
	}
	if len(files) == 0 {
		glog.Exit("No SGF files specified -- must be specified as args to katalyze")
	}

	an := katago.New(model, analysisConfig, *analysisThreads)
	if err = an.Start(); err != nil {
		glog.Exitf("error booting Katago: %v", err)
	}
	var store storage.Filestore
	if storageProvider != nil && *storageProvider == "gcp" {
		glog.Exitf("CloudStorage is not yet implemented")
	} else {
		store, err = getDiskStore()
	}
	if err != nil {
		glog.Exit(err)
	}
	proc := &problemProcessor{
		an: an,
		fs: store,
	}

	if err = proc.genProblems(files); err != nil {
		glog.Exitf("error processing files: %v", err)
	}

	an.Stop()
}

func getDiskStore() (storage.Filestore, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	outDir := cwd
	if *outputDir != "" {
		outDir = path.Join(cwd, *outputDir)
		if _, errz := os.Stat(outDir); os.IsNotExist(errz) {
			os.Mkdir(outDir, storage.DefaultDirPerms)
		} else if errz != nil {
			glog.Exit(errz)
		}
	}
	glog.Infof("Writing files to output dir %v", outDir)
	storage, err := storage.NewDiskStore(outDir)
	if err != nil {
		return nil, err
	}
	return storage, nil
}

func filterSGFs(args []string) ([]string, error) {
	var out []string
	for _, fname := range args {
		fi, err := os.Stat(fname)
		if err != nil {
			return nil, fmt.Errorf("error stating file %v: %v", fname, err)
		}

		var files []string
		switch m := fi.Mode(); {
		case m.IsDir():
			list, err := ioutil.ReadDir(fname)
			if err != nil {
				return nil, fmt.Errorf("error reading dir %v: %v", fname, err)
			}
			for _, l := range list {
				files = append(files, path.Join(fname, l.Name()))
			}
		case m.IsRegular():
			files = []string{fname}
		}
		for _, fin := range files {
			if strings.HasSuffix(fin, ".sgf") {
				out = append(out, fin)
			}
		}
	}
	return out, nil
}
