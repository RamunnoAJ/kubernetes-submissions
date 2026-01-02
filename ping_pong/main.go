package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	var err error
	// Retry connection loop
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

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS pings (id SERIAL PRIMARY KEY, timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func getCount() int {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM pings").Scan(&count)
	if err != nil {
		log.Printf("Failed to get count: %v", err)
		return 0
	}
	return count
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	initDB()

	http.HandleFunc("/pingpong", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pingpong" {
			http.NotFound(w, r)
			return
		}

		_, err := db.Exec("INSERT INTO pings DEFAULT VALUES")
		if err != nil {
			http.Error(w, "Failed to record ping", http.StatusInternalServerError)
			log.Printf("Failed to insert ping: %v", err)
			return
		}

		count := getCount()
		fmt.Fprintf(w, "pong %d\n", count)
	})

	http.HandleFunc("/pings", func(w http.ResponseWriter, r *http.Request) {
		count := getCount()
		fmt.Fprintf(w, "%d", count)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	fmt.Printf("Server started in port %s\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
