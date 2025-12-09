package handler

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"webhook-handler/internal/config"
	"webhook-handler/internal/github"
	"webhook-handler/internal/queue"
)

// WebhookHandler returns an http.HandlerFunc that performs validation, responds, and queues.
func WebhookHandler(cfg *config.Config, publisher *queue.ServiceBusPublisher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// 1. READ RAW BODY
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Handler Error: Failed to read request body: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		r.Body.Close()

		// 2. VALIDATE SIGNATURE
		signature := r.Header.Get("X-Hub-Signature-256")
		if !github.ValidateSignature(signature, body, cfg.GitHubSecret) {
			log.Printf("Handler Error: Signature validation failed for IP %s", r.RemoteAddr)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// 3. EXTRACT FOR LOGGING ONLY
		event := r.Header.Get("X-GitHub-Event")
		// NOTE: Extracting repo name would require JSON unmarshaling, which we avoid
		// to keep the body in its raw state for the queue. Logging the event is sufficient.
		log.Printf("Webhook Received: Event=%s, Length=%d bytes, Signature validated.", event, len(body))


		// 4. IMMEDIATE RESPONSE TO GITHUB (HTTP 200 OK)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Received and queueing for processing."))
		
		// 5. QUEUE THE RAW PAYLOAD ASYNCHRONOUSLY
		// Use a goroutine to ensure queueing does not delay the handler's response.
		// The error handling in the goroutine is crucial.
		go func() {
			// Create a background context for the Service Bus operation
			ctx := context.Background() 
			if err := publisher.Publish(ctx, body); err != nil {
				// LOGGING CRITICAL FAILURE: The webhook was received and validated but failed to queue.
				log.Printf("QUEUE FAILURE: Failed to publish validated payload to Service Bus (Event: %s): %v", event, err)
			} else {
				log.Printf("Queue Success: Payload successfully pushed to Service Bus (Event: %s).", event)
			}
		}()
	}
}
