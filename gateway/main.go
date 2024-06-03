package main

import (
	"log"
	"net/http"

	common "github.com/Far-sa/commons"
	"github.com/Far-sa/gateway/handler"

	_ "github.com/joho/godotenv/autoload"
)

var (
	httpAddr = common.EnvString("HTTP_ADDR", ":3000")
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
