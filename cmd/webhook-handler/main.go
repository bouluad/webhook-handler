package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"webhook-handler/internal/config"
	"webhook-handler/internal/handler"
	"webhook-handler/internal/queue"
)

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Initialize Azure Service Bus Publisher
	publisher, err := queue.NewServiceBusPublisher(cfg.ServiceBusConnectionString, cfg.ServiceBusQueueName)
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize Service Bus publisher: %v", err)
	}
	defer publisher.Close(context.Background())
    
    // 
    // The diagram visually shows the data path: GitHub -> AKS Ingress -> Go Handler -> ASB Queue -> On-prem Jenkins Tool.

	// 3. Setup HTTP Router
	mux := http.NewServeMux()
	mux.Handle("/webhook", handler.WebhookHandler(cfg, publisher))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 4. Start HTTP Server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
		// Recommended timeouts for production-ready servers
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second, 
		IdleTimeout:  120 * time.Second,
	}

	// 5. Graceful Shutdown Implementation (Critical for AKS)
	go func() {
		log.Printf("Server listening on port %s...", cfg.Port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("FATAL: HTTP server ListenAndServe: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop 

	log.Println("Shutting down server...")

	// Shutdown the server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("FATAL: Server shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped.")
}
