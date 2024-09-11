package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gorilla/mux"
)

// Message struct for JSON payload
type Message struct {
	Text string `json:"text"`
}

func main() {
	// Configure Kafka producer
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true // Enable success responses for messages

	// Create a new Kafka producer
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		fmt.Println("Error creating Kafka producer:", err)
		return
	}
	defer producer.Close()

	// Create a new router
	router := mux.NewRouter()

	// Define a route for sending messages
	router.HandleFunc("/send", sendMessage(producer)).Methods("POST")

	// Start the server
	fmt.Println("Server listening on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

// sendMessage function handles the POST request to send messages
func sendMessage(producer sarama.SyncProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the JSON request body
		var message Message
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create a new Kafka message
		kafkaMsg := &sarama.ProducerMessage{
			Topic: "my_topic",
			Value: sarama.StringEncoder(message.Text),
		}

		// Send the message to Kafka
		partition, offset, err := producer.SendMessage(kafkaMsg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("Message sent to topic: %s, partition: %d, offset: %d\n", kafkaMsg.Topic, partition, offset)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Message sent successfully\n")) // Return a success message to the client
	}
}
