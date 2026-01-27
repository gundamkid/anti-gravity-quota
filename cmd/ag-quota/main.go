package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gundamkid/anti-gravity-quota/internal/api"
	"github.com/gundamkid/anti-gravity-quota/internal/auth"
	"github.com/gundamkid/anti-gravity-quota/internal/config"
	"github.com/gundamkid/anti-gravity-quota/internal/ui"
	"github.com/spf13/cobra"
)

var (
	version     = "0.1.0"
	jsonOutput  bool
	accountFlag string
	allFlag     bool
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
		runQuota(cmd, args)
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

// runQuota handles the quota command
func runQuota(cmd *cobra.Command, args []string) {
	// Handle --all flag
	if allFlag {
		runQuotaForAllAccounts()
		return
	}

	// Handle --account flag
	if accountFlag != "" {
		runQuotaForAccount(accountFlag)
		return
	}

	// Default: check quota for default account
	_, err := auth.LoadToken()
	if err != nil {
		if jsonOutput {
			fmt.Fprintf(os.Stderr, `{"error": "not logged in"}%s`, "\n")
		} else {
			ui.DisplayNotLoggedIn()
		}
		os.Exit(1)
	}

	// Show loading message (only if not JSON output)
	if !jsonOutput {
		fmt.Println()
		fmt.Print("Fetching quota information... ")
	}

	// Create API client
	client := api.NewClient()

	// Get quota info
	quotaInfo, err := client.GetQuotaInfo()
	if err != nil {
		if jsonOutput {
			fmt.Fprintf(os.Stderr, `{"error": "failed to fetch quota", "message": "%s"}%s`, err.Error(), "\n")
		} else {
			fmt.Println()
			ui.DisplayError("Failed to fetch quota information", err)
			fmt.Println()
			fmt.Println("Possible issues:")
			fmt.Println("  • Token may have expired (run 'ag-quota login' to re-authenticate)")
			fmt.Println("  • Network connection issues")
			fmt.Println("  • API service may be temporarily unavailable")
		}
		os.Exit(1)
	}

	if !jsonOutput {
		fmt.Println(color.GreenString("✓"))
	}

	// Display quota information
	if jsonOutput {
		if err := ui.DisplayQuotaSummaryJSON(quotaInfo); err != nil {
			fmt.Fprintf(os.Stderr, "Error displaying JSON: %v\n", err)
			os.Exit(1)
		}
	} else {
		ui.DisplayQuotaSummary(quotaInfo)
	}
}

// runQuotaForAccount fetches and displays quota for a specific account
func runQuotaForAccount(email string) {
	if !jsonOutput {
		fmt.Println()
		fmt.Printf("Fetching quota for %s... ", email)
	}

	client := api.NewClient()
	quotaInfo, err := client.GetQuotaInfoForAccount(email)
	if err != nil {
		if jsonOutput {
			fmt.Fprintf(os.Stderr, `{"error": "failed to fetch quota", "account": "%s", "message": "%s"}%s`, email, err.Error(), "\n")
		} else {
			fmt.Println()
			ui.DisplayError(fmt.Sprintf("Failed to fetch quota for %s", email), err)
		}
		os.Exit(1)
	}

	if !jsonOutput {
		fmt.Println(color.GreenString("✓"))
	}

	if jsonOutput {
		if err := ui.DisplayQuotaSummaryJSON(quotaInfo); err != nil {
			fmt.Fprintf(os.Stderr, "Error displaying JSON: %v\n", err)
			os.Exit(1)
		}
	} else {
		ui.DisplayQuotaSummary(quotaInfo)
	}
}

// runQuotaForAllAccounts fetches and displays quota for all saved accounts
func runQuotaForAllAccounts() {
	mgr, err := auth.NewAccountManager()
	if err != nil {
		if jsonOutput {
			fmt.Fprintf(os.Stderr, `{"error": "failed to initialize account manager", "message": "%s"}%s`, err.Error(), "\n")
		} else {
			ui.DisplayError("Failed to initialize account manager", err)
		}
		os.Exit(1)
	}

	accounts, err := mgr.ListAccounts()
	if err != nil {
		if jsonOutput {
			fmt.Fprintf(os.Stderr, `{"error": "failed to list accounts", "message": "%s"}%s`, err.Error(), "\n")
		} else {
			ui.DisplayError("Failed to list accounts", err)
		}
		os.Exit(1)
	}

	if len(accounts) == 0 {
		if jsonOutput {
			fmt.Fprintf(os.Stderr, `{"error": "no accounts found"}%s`, "\n")
		} else {
			color.Yellow("No accounts found. Please run 'ag-quota login' first.")
		}
		os.Exit(1)
	}

	if !jsonOutput {
		fmt.Println()
		fmt.Printf("Fetching quota for %d account(s)...\n", len(accounts))
		fmt.Println()
	}

	// Fetch quota for each account
	client := api.NewClient()
	var quotaResults []*ui.AccountQuotaResult

	for _, acc := range accounts {
		if !jsonOutput {
			fmt.Printf("  • %s... ", acc.Email)
		}

		quotaInfo, err := client.GetQuotaInfoForAccount(acc.Email)
		if err != nil {
			if !jsonOutput {
				fmt.Println(color.RedString("✗"))
				fmt.Printf("    Error: %v\n", err)
			}
			quotaResults = append(quotaResults, &ui.AccountQuotaResult{
				Email: acc.Email,
				Error: err.Error(),
			})
			continue
		}

		if !jsonOutput {
			fmt.Println(color.GreenString("✓"))
		}

		quotaResults = append(quotaResults, &ui.AccountQuotaResult{
			Email:        acc.Email,
			QuotaSummary: quotaInfo,
		})
	}

	if !jsonOutput {
		fmt.Println()
	}

	// Display results
	if jsonOutput {
		if err := ui.DisplayAllAccountsQuotaJSON(quotaResults); err != nil {
			fmt.Fprintf(os.Stderr, "Error displaying JSON: %v\n", err)
			os.Exit(1)
		}
	} else {
		ui.DisplayAllAccountsQuota(quotaResults)
	}
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

	// Add flags
	quotaCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")
	quotaCmd.Flags().StringVar(&accountFlag, "account", "", "Check quota for specific account")
	quotaCmd.Flags().BoolVar(&allFlag, "all", false, "Check quota for all accounts")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
