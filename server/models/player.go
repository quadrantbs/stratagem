package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
	ID           primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	Username     string                 `json:"username" bson:"username"`
	Email        string                 `json:"email" bson:"email"`
	Password     string                 `json:"password" bson:"password"`
	PhotoProfile string                 `json:"photo_profile" bson:"photo_profile"`
	Role         string                 `json:"role" bson:"role"`
	Data         map[string]interface{} `json:"data" bson:"data"`
}
