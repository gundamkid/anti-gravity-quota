package main

import (
	"fmt"
	"os"

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
		fmt.Println("Login command - Coming soon!")
		// TODO: Implement OAuth2 login flow
	},
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",
	Long:  `Display current authentication status and account information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Status command - Coming soon!")
		// TODO: Implement status check
	},
}

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout and clear stored tokens",
	Long:  `Remove stored authentication tokens and logout from the current account.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Logout command - Coming soon!")
		// TODO: Implement logout
	},
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
