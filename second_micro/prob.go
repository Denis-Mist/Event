package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var secretKey = []byte("secret_key") // Replace with your actual secret key

func main() {
	app := fiber.New()
	var err error
	kafkaConsumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaConsumer.Close()

	topicPartitions, err := kafkaConsumer.ConsumePartition("users", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal(err)
	}

	// Define a handler for POST requests to /word
	app.Post("/word", func(c *fiber.Ctx) error {
		// Decode the request body
		var wordRequest struct {
			Word string `json:"word"`
		}

		if err := c.BodyParser(&wordRequest); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		fmt.Printf("Received word: %s\n", wordRequest.Word)

		// You can store the word in the database or do other processing here

		return c.JSON(fiber.Map{"message": "Word received successfully"})
	})

	// Start a separate goroutine to consume messages from Kafka
	go func() {
		for {
			select {
			case msg := <-topicPartitions.Messages():
				fmt.Printf("Received message with value: %s\n", string(msg.Value))

				// Parse the message as JSON
				var messageData map[string]string
				err := json.Unmarshal(msg.Value, &messageData)
				if err != nil {
					log.Println("Error unmarshaling message:", err)
					continue
				}

				// Extract the token from the message data
				tokenString := messageData["token"]

				// Parse the JWT token
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					return secretKey, nil
				})
				if err != nil {
					log.Println("Error parsing token:", err)
					continue
				}
				fmt.Println(token)

				// Make the POST request
				req, err := http.NewRequest("POST", "http://localhost:5050/word", bytes.NewBuffer([]byte(fmt.Sprintf(`{"word": "hello"}`))))
				if err != nil {
					log.Println("Error creating POST request:", err)
					continue
				}
				req.Header.Set("Authorization", "Bearer "+tokenString)

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Println("Error sending POST request:", err)
					continue
				}
				defer resp.Body.Close()

				fmt.Println("POST request sent successfully")
			}
		}
	}()

	log.Fatal(app.Listen(":3000")) // Start the Fiber app on port 3000
}
