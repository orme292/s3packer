package main

import (
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func randomString(n int) string {
	var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, n)
	gen := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range bytes {
		bytes[i] = alphabet[gen.Intn(len(alphabet))]
	}

	return string(bytes)
}

func returnFileNames(path string) (files []string, err error) {
	path, _ = filepath.Abs(path)
	de, err := os.ReadDir(path)

	for _, file := range de {
		if file.IsDir() {
			continue
		}
		files = append(files, file.Name())
	}
	return
}
