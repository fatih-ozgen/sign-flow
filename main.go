package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.Println("Starting application...")

	// Initialize database connection
	log.Println("Initializing database connection...")
	err := initDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection initialized successfully")

	// Set up routes
	log.Println("Setting up routes...")
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/signin", signinHandler)
	http.HandleFunc("/users", getUsersHandler)

	// Add a simple health check route
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Health check requested")
		fmt.Fprintf(w, "Server is up and running")
	})

	log.Println("Routes set up completed")

	// Use http.Server for more control
	server := &http.Server{
		Addr:     ":8080",
		Handler:  nil, // Use default ServeMux
		ErrorLog: log.New(os.Stderr, "HTTP Server Error: ", log.Ldate|log.Ltime|log.Lshortfile),
	}

	// Start server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		log.Println("Attempting to start server on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Wait for the server to start or encounter an error
	select {
	case err := <-serverErr:
		log.Fatalf("Failed to start server: %v", err)
	case <-time.After(2 * time.Second):
		// Check if the server is actually listening
		conn, err := net.DialTimeout("tcp", "localhost:8080", time.Second)
		if err != nil {
			log.Fatalf("Failed to connect to the server: %v", err)
		} else {
			conn.Close()
			log.Println("Successfully connected to the server")
		}
	}

	log.Println("Server is ready. Press Ctrl+C to shut down.")

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Interrupt received, server is shutting down...")
	
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server has been gracefully shut down")
}
