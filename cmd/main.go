package main

import (
	"log"
	"net/http"
	"os"

	"alifattahi.ir/go-dns-resolver/config"
	"alifattahi.ir/go-dns-resolver/handlers"
	"alifattahi.ir/go-dns-resolver/metrics"
	"alifattahi.ir/go-dns-resolver/migrations"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	if err := migrations.CreateTableIfNotExists(db); err != nil {
		log.Fatalf("Failed to create table: %v", err) 
	}
	defer db.Close()

	// Register metrics
	metrics.Register()

	// Initialize HTTP routes
	http.HandleFunc("/resolve", handlers.ResolveDomainHandler(db))
	http.HandleFunc("/ready", handlers.ReadinessHandler(db))
	http.HandleFunc("/healthy", handlers.LivenessHandler)
	http.Handle("/metrics", promhttp.Handler())

	// Start the server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Starting server on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
