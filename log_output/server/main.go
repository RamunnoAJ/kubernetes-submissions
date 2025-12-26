package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	path := os.Getenv("RANDOM_FILE")
	if path == "" {
		path = "/app/files/random.txt"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(path)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot read random.txt: %v", err), http.StatusInternalServerError)
			return
		}

		ppPath := "/app/files/pingpong_count.txt"

		ppData, err := os.ReadFile(ppPath)
		pingPongCount := "0"
		if err == nil {
			pingPongCount = string(ppData)
		}

		fmt.Fprintf(w, "%sPing / Pongs: %s\n", data, pingPongCount)
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
