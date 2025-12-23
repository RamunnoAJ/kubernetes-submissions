package main

import (
	"fmt"
	"math/rand/v2"
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
	randomStr := fmt.Sprintf(
		"%s-%s-%s-%s-%s",
		GenerateRandomString(8),
		GenerateRandomString(4),
		GenerateRandomString(4),
		GenerateRandomString(4),
		GenerateRandomString(12),
	)

	for {
		fmt.Printf("%s: %s\n", time.Now().Format("2006-01-02T15:04:05.705Z"), randomStr)
		time.Sleep(time.Second * 5)
	}
}
