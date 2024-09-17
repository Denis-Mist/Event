package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
)

var DbData = map[string]string{
	"host":     "localhost",
	"port":     "8080", // на 5432
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
	Word  string `json:"word"`
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

	go startServer() // Start the server in a separate goroutine

	for {
		select {
		case msg := <-topicPartitions.Messages():
			fmt.Printf("Received message with value: %s\n", string(msg.Value))
			var token tokenStruct
			err := json.Unmarshal(msg.Value, &token)
			if err != nil {
				log.Println(err)
				fmt.Print("tut :)")
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
			ch := make(chan struct{})
			go func() {
				wordResponse, err := makePostRequest(user, token.Token, token.Word)
				if err != nil {
					log.Println(err)
					ch <- struct{}{}
					return
				}

				// Store word in database
				db, err := createDatabaseIfNotExists()
				if err != nil {
					log.Println(err)
					ch <- struct{}{}
					return
				}
				defer db.Close()

				_, err = db.Exec("INSERT INTO words (user_id, word) VALUES ($1, $2)", user.ID, wordResponse)
				if err != nil {
					log.Println(err)
					ch <- struct{}{}
					return
				}

				ch <- struct{}{}
			}()

			// Wait for the response from the makePostRequest function
			<-ch
		}
	}
}

func makePostRequest(user *User, token string, word string) (string, error) {
	client := &http.Client{}

	// Create a new JSON object with the word field set to the actual word
	wordJSON := []byte(fmt.Sprintf(`{"word": "%s"}`, word))

	// Create a new request with the POST method and the word JSON object as the request body
	req, err := http.NewRequest("POST", "http://localhost:5050/word", bytes.NewBuffer(wordJSON))
	if err != nil {
		return "", err
	}

	// Set the Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	// Set the Authorization header to the JWT token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Return the response body as a string
	return string(body), nil
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

func startServer() {
	http.HandleFunc("/word", func(w http.ResponseWriter, r *http.Request) {
		var wordRequest struct {
			Word string `json:"word"`
		}

		err := json.NewDecoder(r.Body).Decode(&wordRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		word := wordRequest.Word
		json.NewEncoder(w).Encode(map[string]string{"word": word})
	})

	log.Fatal(http.ListenAndServe(":5050", nil))
}
