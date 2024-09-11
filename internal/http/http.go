package http

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

var logsMutex sync.Mutex
var logs []string

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Check status of worker and timer threads
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	logsMutex.Lock()
	defer logsMutex.Unlock()

	json.NewEncoder(w).Encode(logs)
}

func reportToAPI(endpoint string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to send data to API: %v", err)
	}
	defer resp.Body.Close()
}
