package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const imagePath = "/app/images/image.jpg"

func checkAndDownloadImage() {
	// Check if file exists and is recent
	info, err := os.Stat(imagePath)
	if err == nil {
		if time.Since(info.ModTime()) < 10*time.Minute {
			return
		}
	}

	fmt.Println("Downloading new image...")
	resp, err := http.Get("https://picsum.photos/1200")
	if err != nil {
		log.Printf("Failed to fetch image: %v", err)
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(imagePath)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Printf("Failed to save image: %v", err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := os.MkdirAll(filepath.Dir(imagePath), 0o755); err != nil {
		log.Printf("Failed to create image directory: %v", err)
	}

	fs := http.FileServer(http.Dir("./static"))
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Handle path prefix from Ingress
		path := r.URL.Path
		if strings.HasPrefix(path, "/the-project") {
			path = strings.TrimPrefix(path, "/the-project")
			if path == "" {
				path = "/"
			}
		}

		if path == "/image.jpg" {
			checkAndDownloadImage()
			http.ServeFile(w, r, imagePath)
			return
		}
		if path == "/broken" {
			os.Exit(1)
			return
		}
		
		// For serving static files, we need to strip the prefix from the request
		// so the FileServer sees relative paths.
		if strings.HasPrefix(r.URL.Path, "/the-project") {
			http.StripPrefix("/the-project", fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})

	fmt.Printf("Server started in port %s\n", port)
// ...

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
