package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Custom client to NOT follow redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", "https://en.wikipedia.org/wiki/Special:Random", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to fetch random wikipedia article: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Response Status: %s\n", resp.Status)

	finalURL := ""
	if resp.StatusCode == 301 || resp.StatusCode == 302 || resp.StatusCode == 303 || resp.StatusCode == 307 {
		finalURL = resp.Header.Get("Location")
	} else {
		finalURL = resp.Request.URL.String()
	}

	if finalURL == "" {
		log.Fatal("Could not determine final URL")
	}

	fmt.Printf("Found article: %s\n", finalURL)

	todo := map[string]string{
		"text": fmt.Sprintf("Read %s", finalURL),
	}
	
	jsonBody, err := json.Marshal(todo)
	if err != nil {
		log.Fatalf("Failed to marshal json: %v", err)
	}

	backendURL := "http://todo-backend-svc:8080/todos"
	
	postResp, err := http.Post(backendURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatalf("Failed to post todo: %v", err)
	}
	defer postResp.Body.Close()

	if postResp.StatusCode >= 400 {
		log.Fatalf("Backend returned status: %d", postResp.StatusCode)
	}

	fmt.Println("Successfully created todo")
}
