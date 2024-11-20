package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"golang.org/x/oauth2"
)

type Person struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	app := fiber.New()

	// Add logger middleware
	app.Use(logger.New())

	app.Post("/authorize", func(c *fiber.Ctx) error {
		person := new(Person)
		if err := c.BodyParser(person); err != nil {
			return err
		}

		fmt.Printf("Received person information: %+v\n", person)

		// Use the person information to make the OAuth request
		conf := &oauth2.Config{
			ClientID:     "1234567890",
			ClientSecret: "abcdefghijklmnopqrstuvwxyz",
			Scopes:       []string{"email", "profile"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://oauth.vk.com/authorize",
				TokenURL: "https://oauth.vk.com/access_token",
			},
		}

		// Redirect user to consent page to ask for permission
		// for the scopes specified above.
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

		// Use the authorization code that is pushed to the redirect
		// URL. Exchange will do the handshake to retrieve the
		// initial access token. The HTTP Client returned by
		// conf.Client will refresh the token as necessary.
		var code string
		if _, err := fmt.Scan(&code); err != nil {
			log.Fatal(err)
		}

		tok, err := conf.Exchange(context.Background(), code)
		if err != nil {
			log.Fatal(err)
		}

		// Get the user ID from the token
		userID := "unknown" // Replace with actual user ID retrieval logic
		// For example, if you're using VK API, you can use the following code:
		// userID, err := get userIDFromVK(tok.AccessToken)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// Connect to PostgreSQL database
		dbConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			"localhost", "5432", "postgres", "ghbdtn", "users")
		db, err := sql.Open("postgres", dbConn)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Store access token in database
		_, err = db.Exec("INSERT INTO access_tokens (user_id, access_token, expires_at) VALUES ($1, $2, $3)", userID, tok.AccessToken, tok.Expiry)
		if err != nil {
			log.Fatal(err)
		}

		// Create Kafka producer
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		config.Producer.Return.Errors = true

		producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
		if err != nil {
			log.Fatal(err)
		}
		defer producer.Close()

		// Produce a message to a Kafka topic
		msg := &sarama.ProducerMessage{
			Topic: "vk_authorizations",
			Value: sarama.StringEncoder(fmt.Sprintf("User %s authorized with access token %s", userID, tok.AccessToken)),
		}

		_, _, err = producer.SendMessage(msg)
		if err != nil {
			log.Fatal(err)
		}

		return c.JSON(fiber.Map{"message": "Authorization successful"})
	})

	app.Listen(":3000")
}
