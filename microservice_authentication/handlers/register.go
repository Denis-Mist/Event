package handlers

import (
	"encoding/json"
	"net/http"

	"models"
	"utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate user input
	if user.Username == "" || user.Password == "" || user.Email == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Create a new user
	db, err := db.CreateDatabaseIfNotExists()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO users (username, password, email) VALUES ($1, $2, $3)", user.Username, user.Password, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
