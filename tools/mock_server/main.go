package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"
)

type ModelQuota struct {
	Name       string    `json:"name"`
	QuotaUsed  int64     `json:"quotaUsed"`
	QuotaLimit int64     `json:"quotaLimit"`
	ResetAt    time.Time `json:"resetAt"`
}

type UserStatusResponse struct {
	Models []ModelQuota `json:"models"`
}

func main() {
	port := flag.Int("port", 8085, "Port to listen on")
	flag.Parse()

	http.HandleFunc("/GetUserStatus", func(w http.ResponseWriter, r *http.Request) {
		// Verify auth header presence (simple check)
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		resp := UserStatusResponse{
			Models: []ModelQuota{
				{
					Name:       "Gemini 3 Pro",
					QuotaUsed:  1500,
					QuotaLimit: 5000,
					ResetAt:    time.Now().Add(4 * time.Hour),
				},
				{
					Name:       "Claude 3.5 Sonnet",
					QuotaUsed:  200,
					QuotaLimit: 1000,
					ResetAt:    time.Now().Add(4 * time.Hour),
				},
				{
					Name:       "GPT-4o",
					QuotaUsed:  0,
					QuotaLimit: 2000,
					ResetAt:    time.Now().Add(4 * time.Hour),
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Printf("Mock server listening on port %d\n", *port)
	// Keep the process running and clearly named so detector finds it
	// Name must contain "antigravity-language-server"
	// But the binary name is main or mock_server. 
	// My detector checks `cmdline`. 
	// So if I run `go run main.go --port=8085`, the cmdline is `.../exe/main --port=8085` or similar.
	// It might NOT contain "antigravity-language-server".
	// I should build it and name it `antigravity-language-server-mock`.
	
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		panic(err)
	}
}
