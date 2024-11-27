package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

// Database connection
var DB *pgx.Conn

// Redis connection
var RDB *redis.Client

// Context for database operations (background, read-only)
var cntx = context.Background()

// MustGetEnv returns the value of the environment variable key or exits the program
func MustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatal("Environment variable " + key + " not set")
	}
	return value
}

// GetFullUrl returns the full URL for the given alias
func GetFullUrl(alias string) (string, error) {
	var url string
	err := DB.QueryRow(cntx, "SELECT url FROM urls WHERE alias = $1", alias).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func GetRedisCache(alias string) (string, error) {
	return RDB.Get(cntx, alias).Result()
}

func SetRedisCache(alias string, url string) error {
	lifeString := os.Getenv("REDIS_TTL")
	if lifeString == "" {
		// Default to 5 min
		lifeString = "300"
	}
	life, err := strconv.ParseInt(lifeString, 10, 64)
	if err != nil {
		// Default to 5 hour
		life = 300
	}
	return RDB.Set(cntx, alias, url, time.Second*time.Duration(life)).Err()
}

// AliasHandler returns the full URL for the given alias
func AliasHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the alias from the URL path
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		//! TODO: Return 404 Not Found Frontend URL
		http.Error(w, "Alias not provided", http.StatusBadRequest)
		return
	}
	// Check if the alias is in the cache
	cache, err := GetRedisCache(path)
	if strings.Trim(cache, " ") != "" {
		http.Redirect(w, r, cache, http.StatusPermanentRedirect)
		return
	}

	// Fetch the full URL from the database
	url, err := GetFullUrl(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// Set the URL in the cache
	err = SetRedisCache(path, url)
	if err != nil {
		fmt.Printf("Unable to set cache: %v", err)
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
	if err := DB.Ping(cntx); err != nil {
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
	conn, err := pgx.Connect(cntx, DB_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	} else {
		DB = conn
		fmt.Println("Connected to database")
	}
	defer DB.Close(cntx)

	// Connect to Redis
	RDB_URL := MustGetEnv("RDB_URL")
	redisOptions, err := redis.ParseURL(RDB_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse Redis URL: %v\n", err)
		os.Exit(1)
	}
	RDB = redis.NewClient(redisOptions)
	if err := RDB.Ping(cntx).Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to Redis: %v\n", err)
		os.Exit(1)
	}
	defer RDB.Close()

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
