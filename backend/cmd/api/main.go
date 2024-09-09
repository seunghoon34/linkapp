package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/seunghoon34/linkapp/backend/internal/handler"
	"github.com/seunghoon34/linkapp/backend/internal/repository"
	"github.com/seunghoon34/linkapp/backend/internal/service"

	"github.com/joho/godotenv"
)

func loadEnv() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		return
	}

	// Move up two levels to reach the backend directory
	backendDir := filepath.Dir(filepath.Dir(currentDir))
	envPath := filepath.Join(backendDir, ".env")

	err = godotenv.Load(envPath)
	if err != nil {
		log.Println("Error loading .env file:", err)
	}
}

func runLinkExpirationTask(userService *service.UserService) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := userService.ExpireLinks(); err != nil {
				log.Printf("Error expiring links: %v", err)
			}
		}
	}
}

func main() {
	loadEnv()
	// Load environment variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}

	// Set up MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB successfully")

	// Initialize database and collections
	database := client.Database("dating_app")
	// Remove the following line:
	// userCollection := database.Collection("users")

	// Initialize repositories
	userRepo := repository.NewUserRepository(database)
	linkRepo := repository.NewLinkRepository(database)
	chatroomRepo := repository.NewChatroomRepository(database)

	// Initialize services
	userService := service.NewUserService(userRepo, linkRepo, chatroomRepo)

	// Start link expiration goroutine
	go runLinkExpirationTask(userService)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)

	// Set up router
	r := mux.NewRouter()

	// Set up routes
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/users/{id}/matches", userHandler.SearchMatches).Methods("GET")
	r.HandleFunc("/users/{id}/location", userHandler.UpdateLocation).Methods("PUT")
	r.HandleFunc("/users/{id}/start-searching", userHandler.StartSearching).Methods("POST")
	r.HandleFunc("/users/{id}/stop-searching", userHandler.StopSearching).Methods("POST")
	r.HandleFunc("/users/{id}/find-match", userHandler.FindMatch).Methods("GET")
	r.HandleFunc("/users/{userId}/links/{linkId}/respond", userHandler.RespondToLink).Methods("POST")
	r.HandleFunc("/users/{userId}/chatrooms/{chatroomId}/messages", userHandler.SendMessage).Methods("POST")
	r.HandleFunc("/users/{userId}/chatrooms/{chatroomId}/messages", userHandler.GetMessages).Methods("GET")
	r.HandleFunc("/chatrooms/{chatroomId}/unlock", userHandler.UnlockChatroom).Methods("POST")
	r.HandleFunc("/users/{userId}/chatrooms/{chatroomId}/nfc-unlock", userHandler.VerifyNFCAndUnlockChatroom).Methods("POST")

	// Add middleware
	r.Use(loggingMiddleware)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
