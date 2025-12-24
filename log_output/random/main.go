package main

import (
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(n int) string {
	var b strings.Builder
	b.Grow(n)
	for range n {
		b.WriteByte(charset[rand.IntN(len(charset))])
	}
	return b.String()
}

func main() {
	path := os.Getenv("RANDOM_FILE")
	if path == "" {
		path = "/app/files/random.txt"
	}

	interval := 5 * time.Second
	if v := os.Getenv("INTERVAL_SECONDS"); v != "" {
		if s, err := time.ParseDuration(v + "s"); err == nil {
			interval = s
		}
	}

	_ = os.MkdirAll(filepath.Dir(path), 0o755)

	for {
		content := time.Now().Format("2006-01-02T15:04:05.705Z") + ": " + randomString(12) + "\n"
		_ = os.WriteFile(path, []byte(content), 0o644)
		time.Sleep(interval)
	}
}
