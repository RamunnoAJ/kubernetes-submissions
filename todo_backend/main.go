package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

var db *sql.DB

func initDB() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database: %v. Retrying in 2 seconds...", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS todos (id SERIAL PRIMARY KEY, text TEXT)")
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func getTodos(w http.ResponseWriter, r *http.Request) {

rows, err := db.Query("SELECT id, text FROM todos")
	if err != nil {
		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
		log.Printf("Query error: %v", err)
		return
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Text); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		todos = append(todos, t)
	}

	if todos == nil {
		todos = []Todo{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	// Read body to log it, then restore it or decode from buffer
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	
	// Log the raw body or decoded struct
	var t Todo
	if err := json.Unmarshal(body, &t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received todo creation request: %s", t.Text)

	if len(t.Text) > 140 {
		msg := fmt.Sprintf("Todo text too long (max 140 chars). Received %d chars.", len(t.Text))
		log.Printf("REJECTED: %s", msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	err = db.QueryRow("INSERT INTO todos (text) VALUES ($1) RETURNING id", t.Text).Scan(&t.ID)
	if err != nil {
		http.Error(w, "Failed to save todo", http.StatusInternalServerError)
		log.Printf("Insert error: %v", err)
		return
	}

	log.Printf("Successfully created todo with ID: %d", t.ID)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		
		next.ServeHTTP(w, r)
		
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	initDB()

	mux := http.NewServeMux()
	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getTodos(w, r)
			return
		}

		if r.Method == http.MethodPost {
			createTodo(w, r)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Wrap mux with logging middleware
	handler := loggingMiddleware(mux)

	fmt.Printf("Todo Backend started in port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
