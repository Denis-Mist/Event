package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/IBM/sarama"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
)

var DbData = map[string]string{
	"host":     "localhost",
	"port":     "5432",
	"user":     "postgres",
	"password": "ghbdtn",
	"database": "words",
}

type Word struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	Word   string `json:"word"`
}

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type tokenStruct struct {
	Token string `json:"token"`
}

var kafkaConsumer sarama.Consumer

func main() {
	var err error
	kafkaConsumer, err = sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaConsumer.Close()

	topicPartitions, err := kafkaConsumer.ConsumePartition("users", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case msg := <-topicPartitions.Messages():
			fmt.Printf("Received message with value: %s\n", string(msg.Value))
			var token tokenStruct
			err := json.Unmarshal(msg.Value, &token)
			if err != nil {
				log.Println(err)
				continue
			}

			// Verify JWT token
			user, err := VerifyToken(token.Token)
			if err != nil {
				log.Println(err)
				fmt.Print("tut")
				continue
			}

			// Make POST request to store word
			word, err := makePostRequest(user)
			if err != nil {
				log.Println(err)
				continue
			}

			// Store word in database
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
		}
	}
}

func VerifyToken(token string) (*User, error) {
	key, err := generateSecretKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret key: %w", err)
	}

	// Parse the JWT token and claims
	tokenClaims, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := tokenClaims.Claims.(*jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to get claims from token")
	}

	// Check if the "id" claim exists and convert it
	id, ok := (*claims)["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("id claim is missing or invalid")
	}

	// Create user object
	user := &User{
		ID: uint(id),
	}

	return user, nil
}

func makePostRequest(user *User) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8080/word", strings.NewReader(`{"word": ""}`))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var word string
	err = json.NewDecoder(resp.Body).Decode(&word)
	if err != nil {
		return "", err
	}

	return word, nil
}

func createDatabaseIfNotExists() (*sql.DB, error) {
	// Connect to the default database (usually 'postgres')
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		DbData["host"], DbData["port"], DbData["user"], DbData["password"]))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Check if the database exists
	var exists bool
	err = db.QueryRow("SELECT 1 FROM pg_database WHERE datname = $1", DbData["database"]).Scan(&exists)

	if !exists {
		// Create the database
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", DbData["database"]))
		if err != nil {
			return nil, err
		}
	}

	// Close the connection to the default database
	db.Close()

	// Connect to the newly created database
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DbData["host"], DbData["port"], DbData["user"], DbData["password"], DbData["database"]))
	if err != nil {
		return nil, err
	}

	// Create the words table if it does not exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS words (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        word VARCHAR(50) NOT NULL
    )`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func generateSecretKey() (string, error) {
	secret_key := "secret_key"
	return secret_key, nil
}
