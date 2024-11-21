package types

import (
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct {
	ID           primitive.ObjectID     `json:"_id"`
	Username     string                 `json:"username"`
	Email        string                 `json:"email"`
	PhotoProfile string                 `json:"photo_profile"`
	Role         string                 `json:"role"`
	Data         map[string]interface{} `json:"data"`
	jwt.StandardClaims
}

type contextKey string

const ContextKeyPlayer contextKey = "player"

type Response struct {
	Message string `json:"message"`
}
