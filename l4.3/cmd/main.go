package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AugustSerenity/go-contest-L4/l4.3_Events-calendar/internal/handler"
	"github.com/AugustSerenity/go-contest-L4/l4.3_Events-calendar/internal/service"
	"github.com/AugustSerenity/go-contest-L4/l4.3_Events-calendar/internal/storage"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	storage := storage.New()
	srv := service.New(storage)
	h := handler.New(srv)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: h.Route(),
	}

	go func() {
		log.Printf("Starting server on %s...\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	srv.Stop()

	log.Println("Server exited properly")
}
