package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var mu sync.Mutex

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	filePath := os.Getenv("PINGPONG_FILE")
	if filePath == "" {
		filePath = "/app/files/pingpong_count.txt"
	}

	err := os.MkdirAll(filepath.Dir(filePath), 0o755)
	if err != nil {
		log.Printf("Failed to create directory: %v", err)
	}

	count := 0
	data, err := os.ReadFile(filePath)
	if err == nil {
		val, err := strconv.Atoi(string(data))
		if err == nil {
			count = val
		}
	}

	http.HandleFunc("/pingpong", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pingpong" {
			http.NotFound(w, r)
			return
		}

		mu.Lock()
		count++
		err := os.WriteFile(filePath, []byte(strconv.Itoa(count)), 0o644)
		if err != nil {
			log.Printf("Failed to write file: %v", err)
		}
		currentCount := count
		mu.Unlock()

		fmt.Fprintf(
			w,
			"pong %d\n",
			currentCount,
		)
	})

	fmt.Printf("Server started in port %s\n", port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
