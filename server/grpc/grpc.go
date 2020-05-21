// Package grpc contains helpers for initializing gRPC
package grpc

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	log "github.com/golang/glog"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	pb "github.com/otrego/clamshell/server/api"
	"github.com/otrego/clamshell/server/echo"
)

// Options contains options for running gRPC.
type Options struct {
	// Port to listen on.
	Port int
}

// Run starts a gRPC server.
//
// The server will be shutdown when "ctx" is canceled.
//
// This method of serving both gRPC and gRPC-Gateway is inspired by:
// https://github.com/philips/grpc-gateway-example/blob/master/cmd/serve.go
func Run(opts *Options) {
	ctx := context.Background()

	// Listen on addr
	addr := fmt.Sprintf("localhost:%d", opts.Port)

	// Initialize both the http and grpc handlers.
	grpcServer := initGRPCServer()
	mux := http.NewServeMux()

	serveSwaggerDefinition(mux)
	serveGRPCGateway(ctx, mux)

	log.Infof("listening on %s", addr)
	if err := serve(ctx, addr, mux, grpcServer); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func serve(ctx context.Context, addr string, mux *http.ServeMux, grpcServer *grpc.Server) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: handleGRPCAndHTTP(grpcServer, mux),
	}
	return srv.ListenAndServe()
}

func initGRPCServer() *grpc.Server {
	// TODO(kashomon): It's possible we don't need a gRPC server at all and can
	// just use gRPC Gateway directly and skip serving a gRPC endpoint.
	var gopts []grpc.ServerOption
	grpcServer := grpc.NewServer(gopts...)
	pb.RegisterEchoServiceServer(grpcServer, &echo.EchoServer{})
	return grpcServer
}

// serveSwaggerDefinition serves swagger definitions from /swagger.json
func serveSwaggerDefinition(mux *http.ServeMux) {
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		io.Copy(w, strings.NewReader(pb.Swagger))
	})
}

func serveGRPCGateway(ctx context.Context, mux *http.ServeMux) {
	// TODO(kashomon): Add TLS / credentials.
	gwmux := runtime.NewServeMux()
	pb.RegisterEchoServiceHandlerServer(ctx, gwmux, &echo.EchoServer{})
	mux.Handle("/", gwmux)
}

// handleGRPCAndHTTP muxes between gRPC and HTTP based on the header type.
func handleGRPCAndHTTP(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}
