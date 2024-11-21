package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"stratagem-server/db"
	"stratagem-server/models"
	"stratagem-server/types"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	FieldUsername     = "username"
	FieldPhotoProfile = "photo_profile"
)

// PlayerEdit handles player data modification (username and photo profile)
func PlayerEdit(w http.ResponseWriter, r *http.Request) {
	// Extract the username from the URL
	vars := mux.Vars(r)
	username := vars["username"]
	// Verify the authenticated player is the same as the one being edited
	playerFromContext, ok := r.Context().Value(types.ContextKeyPlayer).(models.Player)
	if !ok {
		http.Error(w, "Failed to retrieve player from context", http.StatusUnauthorized)
		return
	}
	if playerFromContext.Username != username {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var player models.Player
	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if player.Username == "" && player.PhotoProfile == "" {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	mongoClient := db.ConnectToMongo()
	collection := mongoClient.Collection("players")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Validate the new username does not exist in the database
	if player.Username != "" {
		var existingPlayer models.Player
		err := collection.FindOne(ctx, bson.M{FieldUsername: player.Username}).Decode(&existingPlayer)
		if err != nil && err != mongo.ErrNoDocuments {
			http.Error(w, "Failed to check existing username", http.StatusInternalServerError)
			return
		}
		if err == nil {
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}
	}

	// Update the player's data
	updateFields := bson.M{}
	if player.Username != "" {
		updateFields[FieldUsername] = player.Username
	}
	if player.PhotoProfile != "" {
		updateFields[FieldPhotoProfile] = player.PhotoProfile
	}

	_, err := collection.UpdateOne(ctx, bson.M{FieldUsername: username}, bson.M{"$set": updateFields})
	if err != nil {
		http.Error(w, "Failed to update player data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.Response{Message: "Player data updated successfully"})
}
