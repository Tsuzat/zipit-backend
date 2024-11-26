package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// Database connection
var DB *pgx.Conn

func MustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatal("Environment variable " + key + " not set")
	}
	return value
}

func GetFullUrl(alias string) (string, error) {
	var url string
	err := DB.QueryRow(context.Background(), "SELECT url FROM urls WHERE alias = $1", alias).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func AliasHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the alias from the URL path
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		//! TODO: Return 404 Not Found Frontend URL
		http.Error(w, "Alias not provided", http.StatusBadRequest)
		return
	}

	url, err := GetFullUrl(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}

// HealthHandler returns the current status of the API
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	isDatabase := false
	if err := DB.Ping(context.Background()); err != nil {
		isDatabase = false
	} else {
		isDatabase = true
	}

	// Create health response
	health := map[string]interface{}{
		"status":   "healthy",
		"database": isDatabase,
	}

	// Set content type and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(health)
}

func main() {
	// Register handlers
	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/", AliasHandler)

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	DB_URL := MustGetEnv("DB_URL")
	conn, err := pgx.Connect(context.Background(), DB_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	} else {
		DB = conn
	}
	defer DB.Close(context.Background())

	// Configure server
	port := 8080
	server := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server
	log.Printf("Server starting on port %d", port)
	log.Fatal(server.ListenAndServe())
}
