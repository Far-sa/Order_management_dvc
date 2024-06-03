package main

import (
	"log"
	"net/http"

	"github.com/Far-sa/gateway/handler"
)

var (
	httpAddr = ":8080"
)

func main() {
	mux := http.NewServeMux()
	handler := handler.New()
	handler.RegisterRoutes(mux)

	log.Printf("starting HTTP server on %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("failed to start server", err)
	}
}
