package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gorilla/mux"
)

// Request struct for JSON payload
type Request struct {
	Text string `json:"text"`
}

// Response struct for JSON payload
type Response struct {
	Text string `json:"text"`
}

func main() {
	// Create a new Kafka producer
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Println("Error creating producer:", err)
		return
	}
	defer producer.Close()

	// Create a new router
	router := mux.NewRouter()

	// Define a route for sending requests
	router.HandleFunc("/send", sendMessage(producer)).Methods("POST")

	// Start the server
	go http.ListenAndServe(":8080", router) // Start the server in a separate goroutine

	// Consume responses from Kafka in the main goroutine
	consumeResponses(producer)
}

// sendMessage function handles the POST request to send messages
func sendMessage(producer sarama.SyncProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the JSON request body
		var request Request
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Send the request to Service B using Kafka
		jsonRequest, err := json.Marshal(request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		kafkaMsg := &sarama.ProducerMessage{
			Topic: "request_topic",
			Value: sarama.ByteEncoder(jsonRequest),
		}
		_, _, err = producer.SendMessage(kafkaMsg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Request sent successfully\n")) // Return a success message to the client
	}
}

// consumeResponses function consumes responses from the response topic
func consumeResponses(producer sarama.SyncProducer) {
	// Create a new Kafka consumer
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Println("Error creating consumer:", err)
		return
	}
	defer consumer.Close()

	// Subscribe to the response topic
	topics := []string{"response_topic"}
	partitions, err := consumer.Partitions(topics[0])
	if err != nil {
		fmt.Println("Error getting partitions:", err)
		return
	}

	// Consume responses from the topic
	for _, partition := range partitions {
		pc, err := consumer.ConsumePartition(topics[0], partition, sarama.OffsetNewest)
		if err != nil {
			fmt.Println("Error consuming partition:", err)
			return
		}
		defer pc.Close()

		for msg := range pc.Messages() {
			var response Response
			err := json.Unmarshal(msg.Value, &response)
			if err != nil {
				fmt.Println("Error unmarshaling response:", err)
				continue // Skip to the next message
			}
			fmt.Printf("Received response: %s\n", response.Text)
		}
	}
}
