package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
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

	// Anti-Gravity OAuth client ID & Secret
	ClientID     = "1071006060591-tmhssin2h21lcre235vtolojh4g403ep.apps.googleusercontent.com"
	ClientSecret = "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"

	// Redirect configuration
	RedirectPort = 42729
	RedirectURI  = "http://127.0.0.1:42729/callback"

	// OAuth2 scopes
	Scopes = "https://www.googleapis.com/auth/cloud-platform https://www.googleapis.com/auth/userinfo.email"
)

// GetOAuthConfig returns the OAuth2 configuration
func GetOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/cloud-platform", "https://www.googleapis.com/auth/userinfo.email"},
		RedirectURL:  RedirectURI,
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

	// Generate state for CSRF protection
	state, err := GenerateState()
	if err != nil {
		return fmt.Errorf("failed to generate state: %w", err)
	}

	// Create a new ServeMux for this login session
	mux := http.NewServeMux()

	// Create channel to receive the result
	resultChan := make(chan LoginResult, 1)

	// Handle OAuth2 callback
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		handleCallback(w, r, oauthConfig, state, "", resultChan)
	})

	// Start local HTTP server for callback on a random port
	// We bind to 127.0.0.1 to avoid firewall prompts on some OSs
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}
	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return fmt.Errorf("failed to get TCP address from listener")
	}
	port := addr.Port

	// Update redirect URL with the actual port
	oauthConfig.RedirectURL = fmt.Sprintf("http://127.0.0.1:%d/callback", port)

	// Start server in background
	server := &http.Server{
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			resultChan <- LoginResult{Error: fmt.Errorf("failed to start callback server: %w", err)}
		}
	}()

	// Build authorization URL
	authURL := oauthConfig.AuthCodeURL(
		state,
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
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Error shutting down server: %v\n", err)
	}

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

	// Exchange code for token
	ctx := context.Background()
	token, err := config.Exchange(ctx, code)
	if err != nil {
		// Log detailed error for debugging
		errMsg := fmt.Sprintf("Failed to exchange token: %v", err)
		fmt.Printf("DEBUG: %s\n", errMsg)
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		resultChan <- LoginResult{Error: fmt.Errorf("failed to exchange token: %w", err)}
		return
	}

	// Fetch user email from Google UserInfo API
	email, err := fetchUserEmail(ctx, token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to fetch user email", http.StatusInternalServerError)
		resultChan <- LoginResult{Error: fmt.Errorf("failed to fetch user email: %w", err)}
		return
	}

	// Create token data
	tokenData := FromOAuth2Token(token, email)

	// Save token
	if err := SaveToken(tokenData); err != nil {
		http.Error(w, "Failed to save token", http.StatusInternalServerError)
		resultChan <- LoginResult{Error: fmt.Errorf("failed to save token: %w", err)}
		return
	}

	// Set the logged-in account as default
	if mgr, err := NewAccountManager(); err == nil {
		_ = mgr.SetDefaultAccount(email)
	}

	// Send success response
	w.Header().Set("Content-Type", "text/html")
	_, _ = fmt.Fprintf(w, `
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

// userInfo represents the response from Google UserInfo API
type userInfo struct {
	Email string `json:"email"`
}

// fetchUserEmail retrieves the user's email address using the access token
func fetchUserEmail(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch user info: status %d", resp.StatusCode)
	}

	var info userInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", err
	}

	return info.Email, nil
}
