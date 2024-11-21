package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"stratagem-server/db"
	"stratagem-server/helpers"
	"stratagem-server/models"
	"stratagem-server/types"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// PlayerLogin handles player login and generates JWT token
func PlayerLogin(w http.ResponseWriter, r *http.Request) {
	var player models.Player

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Connect to MongoDB
	collection := db.ConnectToMongo().Collection("players")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if the player is found by username or email
	var existingPlayer models.Player
	filter := bson.M{}
	if len(player.Username) > 0 {
		filter["username"] = player.Username
	} else if len(player.Email) > 0 {
		filter["email"] = player.Email
	} else {
		http.Error(w, "Username or email is required", http.StatusBadRequest)
		return
	}

	// Search for player by username or email
	err := collection.FindOne(ctx, filter).Decode(&existingPlayer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Username/email or password is incorrect", http.StatusUnauthorized)
			return
		} else {
			http.Error(w, "Error checking username/email", http.StatusInternalServerError)
			return
		}
	}

	// Compare password with hashed password stored in the database
	err = bcrypt.CompareHashAndPassword([]byte(existingPlayer.Password), []byte(player.Password))
	if err != nil {
		http.Error(w, "Username/email or password is incorrect", http.StatusUnauthorized)
		return
	}

	// Create the JWT claims, which includes the username and expiration time
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &types.Claims{
		ID:           existingPlayer.ID,
		Username:     existingPlayer.Username,
		Email:        existingPlayer.Email,
		PhotoProfile: existingPlayer.PhotoProfile,
		Role:         existingPlayer.Role,
		Data:         existingPlayer.Data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	tokenString, err := token.SignedString(helpers.JwtKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return player data (without password) and token
	playerResponse := models.Player{
		ID:           existingPlayer.ID,
		Username:     existingPlayer.Username,
		Email:        existingPlayer.Email,
		PhotoProfile: existingPlayer.PhotoProfile,
		Role:         existingPlayer.Role,
		Data:         existingPlayer.Data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"player":  playerResponse,
		"token":   tokenString,
	})
}
