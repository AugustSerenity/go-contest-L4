package main

import (
	"log"
	"net/http"
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
