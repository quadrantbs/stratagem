package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"stratagem-server/models"
	"strings"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	ID           primitive.ObjectID     `json:"_id"`
	Username     string                 `json:"username"`
	Email        string                 `json:"email"`
	PhotoProfile string                 `json:"photo_profile"`
	Role         string                 `json:"role"`
	Data         map[string]interface{} `json:"data"`
	jwt.StandardClaims
}

// Custom context key type
type contextKey string

// Define specific context key for player
const contextKeyPlayer contextKey = "player"

// AuthMiddleware verifies the JWT token and attaches the user info to the request context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		// Extract token from header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http.Error(w, "Authorization token is missing", http.StatusUnauthorized)
			return
		}

		// Parse and verify the JWT token
		parsedToken, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Return the secret key for validation
			return jwtKey, nil
		})

		if err != nil {
			http.Error(w, "Failed to parse token", http.StatusUnauthorized)
			return
		}

		// Verify if the token is valid
		if claims, ok := parsedToken.Claims.(*Claims); ok && parsedToken.Valid {

			// Create player from claims and add to context
			player := models.Player{
				ID:           claims.ID,
				Username:     claims.Username,
				Email:        claims.Email,
				PhotoProfile: claims.PhotoProfile,
				Role:         claims.Role,
				Data:         claims.Data,
			}

			// Add player to context for downstream handlers
			ctx := context.WithValue(r.Context(), contextKeyPlayer, player)
			fmt.Println("auth middleware ctx:", ctx.Value(contextKeyPlayer))
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
	})
}
