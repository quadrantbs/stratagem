package helpers

import "os"

var JwtKey = []byte(os.Getenv("JWT_SECRET"))
