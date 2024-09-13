package db

import (
	"database/sql"
	"fmt"

	"config"
)

func createDatabaseIfNotExists() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.ConfigDbData["host"], config.ConfigDbData["port"], config.ConfigDbData["user"], config.ConfigDbData["password"], config.ConfigDbData["database"]))
	if err != nil {
		return nil, err
	}

	// Check if the database exists
	var exists bool
	err = db.QueryRow("SELECT 1 FROM pg_database WHERE datname = $1", config.ConfigDbData["database"]).Scan(&exists)
	if err != nil {
		return nil, err
	}

	if !exists {
		// Create the database
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", config.ConfigDbData["database"]))
		if err != nil {
			return nil, err
		}
	}

	// Create the users table if it does not exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) NOT NULL,
        password VARCHAR(100) NOT NULL,
        email VARCHAR(100) NOT NULL
    )`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
