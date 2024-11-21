package middlewares

import (
	"context"
	"net/http"
	"stratagem-server/helpers"
	"stratagem-server/models"
	"stratagem-server/types"
	"strings"

	"github.com/golang-jwt/jwt"
)

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
		parsedToken, err := jwt.ParseWithClaims(tokenString, &types.Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Return the secret key for validation
			return helpers.JwtKey, nil
		})

		if err != nil {
			http.Error(w, "Failed to parse token", http.StatusUnauthorized)
			return
		}

		// Verify if the token is valid
		if claims, ok := parsedToken.Claims.(*types.Claims); ok && parsedToken.Valid {

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
			ctx := context.WithValue(r.Context(), types.ContextKeyPlayer, player)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
	})
}
