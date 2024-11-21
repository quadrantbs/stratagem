package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"stratagem-server/db"
	"stratagem-server/helpers"
	"stratagem-server/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// PlayerRegister handles player registration
func PlayerRegister(w http.ResponseWriter, r *http.Request) {
	var player models.Player

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate email format
	if !helpers.IsValidEmail(player.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Validate username length
	if len(player.Username) < 4 {
		http.Error(w, "Username must be at least 4 characters", http.StatusBadRequest)
		return
	}

	// Validate password length
	if len(player.Password) < 6 {
		http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	// Set default role if not provided
	if player.Role == "" {
		player.Role = "guest"
	}

	// Set default photo profile
	player.PhotoProfile = "/images/default_profile.png"

	// Set default data
	player.Data = map[string]interface{}{}

	// Check if the username or email is already taken
	collection := db.ConnectToMongo().Collection("players")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if the username already exists
	var existingPlayer models.Player
	err := collection.FindOne(ctx, bson.M{"username": player.Username}).Decode(&existingPlayer)
	if err == nil {
		http.Error(w, "Username is already taken", http.StatusConflict)
		return
	}

	// Check if the email already exists
	err = collection.FindOne(ctx, bson.M{"email": player.Email}).Decode(&existingPlayer)
	if err == nil {
		http.Error(w, "Email is already registered", http.StatusConflict)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(player.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	player.Password = string(hashedPassword)

	// Save player to the database
	_, err = collection.InsertOne(ctx, player)
	if err != nil {
		http.Error(w, "Failed to register player", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Player registered successfully"})
}
