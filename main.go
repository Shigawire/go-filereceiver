package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// StoragePath holds the absolute storage path for all uploads
const StoragePath = "/app/storage"

// BasicAuthCredentials holds the servers required basic auth credentials
type BasicAuthCredentials struct {
	username string
	password string
}

var permissableContentTypes = map[string]bool{
	"application/x-gzip": true,
	"application/pdf":    true,
}

func main() {

	basicAuthUsername, ok := os.LookupEnv("BASIC_AUTH_USERNAME")
	if !ok {
		log.Fatal("BASIC_AUTH_USERNAME must be set")
	}

	basicAuthPassword, ok := os.LookupEnv("BASIC_AUTH_PASSWORD")
	if !ok {
		log.Fatal("BASIC_AUTH_PASSWORD must be set")
	}

	creds := BasicAuthCredentials{basicAuthUsername, basicAuthPassword}

	fmt.Print("Serving file writer.\n")

	http.HandleFunc("/receive", creds.receiveHandler)
	http.HandleFunc("/.well-known/health-check", healthCheckHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (creds *BasicAuthCredentials) receiveHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	username, password, ok := r.BasicAuth()

	if !ok {
		fmt.Println("Could not check basic auth credentials.")
		return
	}

	if username != creds.username && password != creds.password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	buffer := make([]byte, 512)
	if _, err := r.Body.Read(buffer); err != nil && err != io.EOF {
		log.Fatal(err)
		return
	}

	contentType := http.DetectContentType(buffer)

	if !permissableContentTypes[contentType] {
		fmt.Printf("Expected ContentType to be any of %v. Got %s.\n", permissableContentTypes, contentType)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	t := time.Now()
	timestampedFilename := t.Format("2006-01-02-15-04-05.pdf")

	filePath := filepath.Join(StoragePath, timestampedFilename)

	file, err := os.Create(filePath)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer file.Close()

	// write the first 512 bytes from buffer to file
	file.Write(buffer)

	// write the rest
	fileWriter := bufio.NewWriter(file)

	if _, err := io.CopyBuffer(fileWriter, r.Body, nil); err != nil {
		log.Fatal(err)
		return
	}

	if err := fileWriter.Flush(); err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("Written file to %s\n", filePath)

	fmt.Fprintf(w, "ok: %s", timestampedFilename)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
