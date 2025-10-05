package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"task-project/db"
	"task-project/routes"

	"github.com/rs/cors"
)

func main() {
	log.Println("Go Running...")

	// err := godotenv.Load(".env.local") //-- try in local
	// if err != nil {
	// 	log.Println("Warning: .env not found, using system env")
	// }
	// Load port from environment
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	// Initialize MongoDB
	mongoStore, err := db.InitMongo()
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}
	defer mongoStore.Client.Disconnect(context.Background())

	// Setup routes
	route := routes.SetupRouter(mongoStore)

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"https://task-project-fe-production.up.railway.app",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(route)

	// Setup HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	// Run server in goroutine
	go func() {
		log.Printf("Server running on port %s...", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
