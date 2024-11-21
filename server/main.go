package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"stratagem-server/controllers"
	"stratagem-server/db"
	"stratagem-server/middlewares"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	database := db.ConnectToMongo()

	fmt.Println("Connected to database:", database.Name())

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Stratagame Backend")
	}).Methods("GET")

	r.HandleFunc("/players/register", controllers.PlayerRegister).Methods("POST")
	r.HandleFunc("/players/login", controllers.PlayerLogin).Methods("POST")
	r.Handle("/players/{username}/edit", middlewares.AuthMiddleware(http.HandlerFunc(controllers.PlayerEdit))).Methods("PUT")

	r.Use(middlewares.EnableCORS)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server started on http://localhost:" + port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
