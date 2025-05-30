package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const port = "8080"
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

type apiConfig struct {
	fileserverHits atomic.Int32
}
