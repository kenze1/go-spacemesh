package api

import (
	"flag"
	"github.com/UnrulyOS/go-unruly/log"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net/http"

	config "github.com/UnrulyOS/go-unruly/app/config"

	gw "github.com/UnrulyOS/go-unruly/api/pb"
)

// A json http server providing the Unruly API.
// Implemented as a grpc gateway. See https://github.com/grpc-ecosystem/grpc-gateway

func run(config *config.Config) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	echoEndpoint := flag.String("api_endpoint", "localhost:"+string(config.GrpcServerPort), "endpoint of api grpc service")

	if err := gw.RegisterUnrulyServiceHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts); err != nil {
		log.Error("Failed to register http endpoint with grpc: %v", err)
	}

	addr := ":" + string(config.JsonServerPort)

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Error("Failed to listen and serve: v%", err)
	}

	return err
}

func StartJsonServer(config *config.Config) {
	if err := run(config); err != nil {
		log.Error("Json serving error: %v", err)
	}
}