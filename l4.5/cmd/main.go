package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"l2.18/internal/handler"
	"l2.18/internal/service"
	"l2.18/internal/storage"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port

	go func() {
		log.Println("Starting pprof on :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	storage := storage.New()
	srv := service.New(storage)
	h := handler.New(srv)

	s := http.Server{
		Addr:    addr,
		Handler: h.Route(),
	}

	log.Printf("Starting server on %s...\n", addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
