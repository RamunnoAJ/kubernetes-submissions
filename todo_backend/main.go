package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	var t Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(t.Text) > 140 {
		http.Error(w, "Todo text too long (max 140 chars)", http.StatusBadRequest)
		return
	}

	err := db.QueryRow("INSERT INTO todos (text) VALUES ($1) RETURNING id", t.Text).Scan(&t.ID)
	if err != nil {
		http.Error(w, "Failed to save todo", http.StatusInternalServerError)
		log.Printf("Insert error: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	initDB()

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
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

	fmt.Printf("Todo Backend started in port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}