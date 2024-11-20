package routes

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate user input
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Find user by username and password
	db, err := gorm.Open("mysql", "user:password@tcp(localhost:3306)/database")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var foundUser User
	db.Where("username = ? AND password = ?", user.Username, user.Password).First(&foundUser)

	if foundUser.ID == 0 {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := GenerateToken(&foundUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
