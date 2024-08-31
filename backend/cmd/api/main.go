package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/seunghoon34/linkapp/backend/internal/handler"
	"github.com/seunghoon34/linkapp/backend/internal/repository"
	"github.com/seunghoon34/linkapp/backend/internal/service"
	"github.com/seunghoon34/linkapp/backend/pkg/db"
)

func main() {
	// Initialize database connection
	database, err := db.NewPostgresConnection("localhost", "5432", "username", "password", "dating_app")
	if err != nil {
		log.Fatalf("Could not initialize database connection: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(database)

	// Initialize services
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")

	// Start server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
