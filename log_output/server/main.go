package main

import (
	"fmt"
	"io"
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

		// Read ConfigMap file
		infoPath := "/app/config/information.txt"
		infoData, err := os.ReadFile(infoPath)
		fileContent := "failed to read file"
		if err == nil {
			fileContent = string(infoData)
		}

		// Read Env Variable
		message := os.Getenv("MESSAGE")

		// Fetch pingpong count
		pingPongCount := "0"
		resp, err := http.Get("http://ping-pong-svc:8080/pings")
		if err == nil {
			body, err := io.ReadAll(resp.Body)
			if err == nil {
				pingPongCount = string(body)
			}
			resp.Body.Close()
		} else {
			log.Printf("Failed to call ping-pong svc: %v", err)
		}

		fmt.Fprintf(w, "file content: %s\n", fileContent)
		fmt.Fprintf(w, "env variable: MESSAGE=%s\n", message)
		fmt.Fprintf(w, "%sPing / Pongs: %s\n", data, pingPongCount)
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
