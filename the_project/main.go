package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const imagePath = "/app/images/image.jpg"

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

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

	tmpl := template.Must(template.ParseFiles("./static/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

		if path == "/todos" && r.Method == http.MethodPost {
			todoText := r.FormValue("todo")
			if todoText != "" {
				todo := Todo{Text: todoText}
				jsonBody, _ := json.Marshal(todo)
				_, err := http.Post("http://todo-backend-svc:8080/todos", "application/json", bytes.NewBuffer(jsonBody))
				if err != nil {
					log.Printf("Failed to create todo: %v", err)
				}
			}
			http.Redirect(w, r, "/the-project/", http.StatusSeeOther)
			return
		}

		resp, err := http.Get("http://todo-backend-svc:8080/todos")
		var todos []Todo
		if err == nil {
			json.NewDecoder(resp.Body).Decode(&todos)
			resp.Body.Close()
		} else {
			log.Printf("Failed to fetch todos: %v", err)
		}

		w.Header().Set("Content-Type", "text/html")
		tmpl.Execute(w, struct{ Todos []Todo }{Todos: todos})
	})

	fmt.Printf("Server started in port %s\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
