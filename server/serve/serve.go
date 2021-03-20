package serve

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	log "github.com/golang/glog"
	"github.com/otrego/clamshell/server/assets"
	"github.com/otrego/clamshell/server/config"
)

func serveOneFile(f []byte) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		io.Copy(w, bytes.NewReader(f))
	}
}

// Run creates an HTTP server.
func Run(opts *config.Spec) {
	http.HandleFunc("/", serveOneFile(assets.Index))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(assets.Assets))))

	addr := fmt.Sprintf("0.0.0.0:%d", opts.Port)

	log.Fatal(http.ListenAndServe(addr, nil))
}
