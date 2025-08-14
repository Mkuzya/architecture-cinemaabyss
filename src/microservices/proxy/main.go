// ./src/microservices/proxy/main.go
package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"
)

var (
	monolithURL     string
	moviesServiceURL string
	gradualMigration bool
	migrationPercent int
)

func init() {
	rand.Seed(time.Now().UnixNano())
	
	// Загрузка конфигурации из переменных окружения
	monolithURL = os.Getenv("MONOLITH_URL")
	moviesServiceURL = os.Getenv("MOVIES_SERVICE_URL")
	
	if os.Getenv("GRADUAL_MIGRATION") == "true" {
		gradualMigration = true
	} else {
		gradualMigration = false
	}
	
	percent, err := strconv.Atoi(os.Getenv("MOVIES_MIGRATION_PERCENT"))
	if err != nil {
		log.Fatalf("Invalid MOVIES_MIGRATION_PERCENT: %v", err)
	}
	migrationPercent = percent
}

func createReverseProxy(target string) *httputil.ReverseProxy {
	url, _ := url.Parse(target)
	return httputil.NewSingleHostReverseProxy(url)
}

func shouldRouteToNewService() bool {
	if !gradualMigration {
		return false
	}
	if migrationPercent >= 100 {
		return true
	}
	if migrationPercent <= 0 {
		return false
	}
	return rand.Intn(100) < migrationPercent
}

func moviesHandler(w http.ResponseWriter, r *http.Request) {
	var target string
	if shouldRouteToNewService() {
		log.Println("Routing to NEW movies service")
		target = moviesServiceURL
	} else {
		log.Println("Routing to LEGACY monolith")
		target = monolithURL
	}
	
	proxy := createReverseProxy(target)
	proxy.ServeHTTP(w, r)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	http.HandleFunc("/api/movies", moviesHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("Starting proxy server on port %s\n", port)
	log.Printf("Migration settings: gradual=%t, percent=%d%%\n", gradualMigration, migrationPercent)
	log.Printf("Monolith URL: %s\n", monolithURL)
	log.Printf("Movies service URL: %s\n", moviesServiceURL)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
