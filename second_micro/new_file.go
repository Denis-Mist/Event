package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	_ "github.com/lib/pq"
)

// ... (Your struct definitions for User, Word, etc.)

var (
	DbData = map[string]string{
		"host":     "localhost",
		"port":     "5432",
		"user":     "postgres",
		"password": "ghbdtn",
		"database": "words",
	}
	kafkaConsumer sarama.Consumer
	secretKey     = "secret_key" // Replace with your actual secret key
)

func main() {
	// ... (Your database setup code)

	var err error
	kafkaConsumer, err = sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaConsumer.Close()

	topicPartitions, err := kafkaConsumer.ConsumePartition("jwt_tokens", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal(err)
	}

	// Start the HTTP server for receiving POST requests
	go startServer()

	for {
		select {
		case msg := <-topicPartitions.Messages():
			fmt.Printf("Received message with value: %s\n", string(msg.Value))

			// Decode the JWT token from the message
			var token string
			err := json.Unmarshal(msg.Value, &token)
			if err != nil {
				log.Println(err)
				continue
			}

			// Verify the JWT token
			user, err := VerifyToken(token)
			if err != nil {
				log.Println(err)
				continue
			}

			// Create a channel to receive the word from the POST request
			wordCh := make(chan string)

			// Start a goroutine to wait for the POST request
			go func(user *User) {
				// Wait for the POST request with a timeout
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				word, err := waitForPostRequest(ctx, user.ID)
				if err != nil {
					log.Printf("Error waiting for POST request: %s", err)
					wordCh <- ""
					return
				}
				wordCh <- word
			}(user)

			// Wait for the word from the POST request or timeout
			select {
			case word := <-wordCh:
				if word == "" {
					// Timeout occurred
					log.Printf("Timeout waiting for POST request for user %d\n", user.ID)
					continue
				}

				// Store the word in the database
				db, err := createDatabaseIfNotExists()
				if err != nil {
					log.Println(err)
					continue
				}
				defer db.Close()

				_, err = db.Exec("INSERT INTO words (user_id, word) VALUES ($1, $2)", user.ID, word)
				if err != nil {
					log.Println(err)
					continue
				}
			case <-ctx.Done():
				// Timeout occurred
				log.Printf("Timeout waiting for POST request for user %d\n", user.ID)
			}
		}
	}
}

// ... (Your VerifyToken and createDatabaseIfNotExists functions)

// Function to handle POST requests
func handleWord(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the request body
	var wordRequest struct {
		Word string `json:"word"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&wordRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Extract the user ID from the context
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	// Send the word to the channel for the corresponding user
	wordCh, ok := r.Context().Value("wordCh").(chan string)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	wordCh <- wordRequest.Word

	// Send a success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Word received successfully"})
}

// Function to wait for a POST request for the given user ID
func waitForPostRequest(ctx context.Context, userID uint) (string, error) {
	// Create a channel to receive the word
	wordCh := make(chan string)

	// Create a request handler that stores the word in the channel
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		// ... (same code as handleWord to decode and store the word in the channel)
	})

	// Create a server with the request handler and set the user ID and word channel in the context
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", userID), // Use the user ID as the port number
		Handler: handler,
	}

	// Create a new context with the user ID and word channel
	requestCtx := context.WithValue(ctx, "userID", userID)
	requestCtx = context.WithValue(requestCtx, "wordCh", wordCh)

	// Start the server in a separate goroutine
	go func() {
		if err := server.Serve(requestCtx); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}()

	// Wait for the word from the channel or the timeout
	select {
	case word := <-wordCh:
		return word, nil
	case <-ctx.Done():
		server.Shutdown(context.Background())
		return "", fmt.Errorf("timed out waiting for POST request")
	}
}

// Function to start the HTTP server
func startServer() {
	http.HandleFunc("/word", handleWord)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
