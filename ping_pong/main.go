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

	count := 0
	http.HandleFunc("/pingpong", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(
			w,
			"pong %d\n",
			count,
		)
		count++
	})

	fmt.Printf("Server started in port %s\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
