package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var mu sync.Mutex

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	count := 0

	http.HandleFunc("/pingpong", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pingpong" {
			http.NotFound(w, r)
			return
		}

		mu.Lock()
		count++
		currentCount := count
		mu.Unlock()

		fmt.Fprintf(
			w,
			"pong %d\n",
			currentCount,
		)
	})

	http.HandleFunc("/pings", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		currentCount := count
		mu.Unlock()

		fmt.Fprintf(w, "%d", currentCount)
	})

	fmt.Printf("Server started in port %s\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
