package main

import (
	"log"
	"net/http"

	common "github.com/Far-sa/commons"
	"github.com/Far-sa/commons/discovery/consul"
	"github.com/Far-sa/gateway/gateway"
	"github.com/Far-sa/gateway/handler"

	_ "github.com/joho/godotenv/autoload"
)

var (
	serviceName = "gateway"
	httpAddr    = common.EnvString("HTTP_ADDR", ":3000")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	gateway.NewGRPCGateway(registry)

	mux := http.NewServeMux()
	handler := handler.New()
	handler.RegisterRoutes(mux)

	log.Printf("starting HTTP server on %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("failed to start server", err)
	}
}
