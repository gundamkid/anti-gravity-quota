package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gundamkid/anti-gravity-quota/internal/auth"
	"github.com/gundamkid/anti-gravity-quota/internal/config"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ag-quota",
	Short: "Anti-Gravity Quota CLI",
	Long: `Anti-Gravity Quota CLI - Check quota for Claude and Gemini models.

This tool allows you to monitor your AI model quota usage through
the Google Cloud Code API.`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		// Default action is to show quota
		quotaCmd.Run(cmd, args)
	},
}

// quotaCmd represents the quota command
var quotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Check quota for all models",
	Long:  `Display quota information for all available AI models (Claude and Gemini).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Quota command - Coming soon!")
		// TODO: Implement quota display
	},
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login with Google account",
	Long:  `Start OAuth2 login flow to authenticate with your Google account.`,
	Run: func(cmd *cobra.Command, args []string) {
		runLogin(cmd, args)
	},
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",
	Long:  `Display current authentication status and account information.`,
	Run: func(cmd *cobra.Command, args []string) {
		runStatus(cmd, args)
	},
}

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout and clear stored tokens",
	Long:  `Remove stored authentication tokens and logout from the current account.`,
	Run: func(cmd *cobra.Command, args []string) {
		runLogout(cmd, args)
	},
}

// runLogin handles the login command
func runLogin(cmd *cobra.Command, args []string) {
	fmt.Println("Starting authentication flow...")
	fmt.Println()

	if err := auth.Login(); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}

// runStatus handles the status command
func runStatus(cmd *cobra.Command, args []string) {
	token, err := auth.LoadToken()
	if err != nil {
		color.Red("Not logged in")
		fmt.Println("Run 'ag-quota login' to authenticate")
		os.Exit(1)
	}

	fmt.Println("Authentication Status")
	fmt.Println("====================")
	fmt.Println()

	if token.Email != "" {
		color.Green("✓ Logged in as: %s", token.Email)
	} else {
		color.Green("✓ Logged in")
	}

	// Check token validity
	if token.IsExpired() {
		color.Yellow("⚠ Token expired")
		if token.RefreshToken != "" {
			fmt.Println("  Token will be automatically refreshed on next use")
		} else {
			fmt.Println("  Please run 'ag-quota login' to re-authenticate")
		}
	} else {
		timeUntilExpiry := time.Until(token.Expiry)
		color.Green("✓ Token valid for: %s", timeUntilExpiry.Round(time.Minute))
	}

	// Show config directory
	configDir, err := config.GetConfigDir()
	if err == nil {
		fmt.Println()
		fmt.Printf("Config directory: %s\n", configDir)
	}
}

// runLogout handles the logout command
func runLogout(cmd *cobra.Command, args []string) {
	// Check if logged in first
	_, err := auth.LoadToken()
	if err != nil {
		color.Yellow("Not logged in")
		return
	}

	// Delete token
	if err := auth.DeleteToken(); err != nil {
		color.Red("Error logging out: %v", err)
		os.Exit(1)
	}

	color.Green("✓ Logged out successfully")
}

func init() {
	// Add subcommands to root
	rootCmd.AddCommand(quotaCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(logoutCmd)

	// Add flags if needed
	// quotaCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
