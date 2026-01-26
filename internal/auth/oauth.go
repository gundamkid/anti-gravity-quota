package auth

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	// Google OAuth2 endpoints
	GoogleAuthURL  = "https://accounts.google.com/o/oauth2/v2/auth"
	GoogleTokenURL = "https://oauth2.googleapis.com/token"

	// Google Cloud Code OAuth client ID (public, from Cloud Code extension)
	ClientID = "764086051850-6qr4p6gpi6hn506pt8ejuq83di341hur.apps.googleusercontent.com"

	// Redirect configuration
	RedirectPort = 8085
	RedirectURI  = "http://localhost:8085/callback"

	// OAuth2 scopes
	Scopes = "openid email profile https://www.googleapis.com/auth/cloud-platform"
)

// GetOAuthConfig returns the OAuth2 configuration
func GetOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID: ClientID,
		Endpoint: google.Endpoint,
		Scopes:   []string{"openid", "email", "profile", "https://www.googleapis.com/auth/cloud-platform"},
		RedirectURL: RedirectURI,
	}
}

// LoginResult contains the result of a login attempt
type LoginResult struct {
	Token *TokenData
	Email string
	Error error
}

// Login initiates the OAuth2 login flow
func Login() error {
	oauthConfig := GetOAuthConfig()

	// Generate PKCE parameters
	verifier, err := GenerateCodeVerifier()
	if err != nil {
		return fmt.Errorf("failed to generate code verifier: %w", err)
	}

	challenge := GenerateCodeChallenge(verifier)

	// Generate state for CSRF protection
	state, err := GenerateState()
	if err != nil {
		return fmt.Errorf("failed to generate state: %w", err)
	}

	// Create channel to receive the result
	resultChan := make(chan LoginResult, 1)

	// Start local HTTP server for callback
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", RedirectPort),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Handle OAuth2 callback
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		handleCallback(w, r, oauthConfig, state, verifier, resultChan)
	})

	// Start server in background
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			resultChan <- LoginResult{Error: fmt.Errorf("failed to start callback server: %w", err)}
		}
	}()

	// Wait a moment for server to start
	time.Sleep(100 * time.Millisecond)

	// Build authorization URL with PKCE
	authURL := oauthConfig.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("access_type", "offline"),
		oauth2.SetAuthURLParam("prompt", "consent"),
	)

	// Open browser
	fmt.Println("Opening browser for authentication...")
	fmt.Println("If browser doesn't open, visit this URL:")
	fmt.Println(authURL)
	fmt.Println()

	if err := openBrowser(authURL); err != nil {
		fmt.Printf("Could not open browser automatically: %v\n", err)
	}

	// Wait for callback result with timeout
	var result LoginResult
	select {
	case result = <-resultChan:
		// Got result
	case <-time.After(5 * time.Minute):
		result = LoginResult{Error: fmt.Errorf("authentication timeout")}
	}

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	if result.Error != nil {
		return result.Error
	}

	fmt.Println("\nLogin successful!")
	if result.Email != "" {
		fmt.Printf("Logged in as: %s\n", result.Email)
	}

	return nil
}

// handleCallback processes the OAuth2 callback
func handleCallback(w http.ResponseWriter, r *http.Request, config *oauth2.Config, expectedState, verifier string, resultChan chan LoginResult) {
	// Check for error in callback
	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		errDesc := r.URL.Query().Get("error_description")
		http.Error(w, "Authentication failed", http.StatusBadRequest)
		resultChan <- LoginResult{Error: fmt.Errorf("authentication failed: %s - %s", errMsg, errDesc)}
		return
	}

	// Verify state parameter (CSRF protection)
	state := r.URL.Query().Get("state")
	if state != expectedState {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		resultChan <- LoginResult{Error: fmt.Errorf("invalid state parameter")}
		return
	}

	// Get authorization code
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No authorization code received", http.StatusBadRequest)
		resultChan <- LoginResult{Error: fmt.Errorf("no authorization code received")}
		return
	}

	// Exchange code for token with PKCE verifier
	ctx := context.Background()
	token, err := config.Exchange(
		ctx,
		code,
		oauth2.SetAuthURLParam("code_verifier", verifier),
	)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		resultChan <- LoginResult{Error: fmt.Errorf("failed to exchange token: %w", err)}
		return
	}

	// Get user email from ID token
	email := ""
	if idToken, ok := token.Extra("id_token").(string); ok {
		// For now, we'll skip parsing the JWT
		// In production, you'd decode the JWT to get email
		_ = idToken
	}

	// Create token data
	tokenData := FromOAuth2Token(token, email)

	// Save token
	if err := SaveToken(tokenData); err != nil {
		http.Error(w, "Failed to save token", http.StatusInternalServerError)
		resultChan <- LoginResult{Error: fmt.Errorf("failed to save token: %w", err)}
		return
	}

	// Send success response
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Authentication Successful</title>
			<style>
				body {
					font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
					display: flex;
					justify-content: center;
					align-items: center;
					height: 100vh;
					margin: 0;
					background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
				}
				.container {
					background: white;
					padding: 3rem;
					border-radius: 10px;
					box-shadow: 0 10px 40px rgba(0,0,0,0.2);
					text-align: center;
					max-width: 500px;
				}
				h1 { color: #333; margin-top: 0; }
				p { color: #666; font-size: 1.1rem; }
				.success { color: #22c55e; font-size: 3rem; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="success">✓</div>
				<h1>Authentication Successful!</h1>
				<p>You have been successfully authenticated.</p>
				<p>You can close this window and return to the terminal.</p>
			</div>
		</body>
		</html>
	`)

	// Send result
	resultChan <- LoginResult{
		Token: tokenData,
		Email: email,
		Error: nil,
	}
}

// openBrowser opens the default browser to the specified URL
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
}
