package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(length int) string {
	var sb strings.Builder
	sb.Grow(length)
	for range length {
		sb.WriteByte(charset[rand.IntN(len(charset))])
	}
	return sb.String()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	randomStr := fmt.Sprintf(
		"%s-%s-%s-%s-%s",
		GenerateRandomString(8),
		GenerateRandomString(4),
		GenerateRandomString(4),
		GenerateRandomString(4),
		GenerateRandomString(12),
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(
			w,
			"%s: %s\n",
			time.Now().Format("2006-01-02T15:04:05.705Z"),
			randomStr,
		)
	})

	fmt.Printf("Server started in port %s\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
