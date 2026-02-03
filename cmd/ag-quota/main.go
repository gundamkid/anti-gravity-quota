package main

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/gundamkid/anti-gravity-quota/internal/api"
	"github.com/gundamkid/anti-gravity-quota/internal/auth"
	"github.com/gundamkid/anti-gravity-quota/internal/config"
	"github.com/gundamkid/anti-gravity-quota/internal/ui"
	"github.com/spf13/cobra"
)

var (
	version       = "0.1.1"
	jsonOutput    bool
	accountFlag   string
	allFlag       bool
	watchInterval int
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
	// Handle watch mode
	if cmd.Flags().Changed("watch") {
		if jsonOutput {
			ui.DisplayError("Flag conflict", fmt.Errorf("--watch cannot be used with --json"))
			os.Exit(1)
		}

		if watchInterval == 0 {
			// If flag is present but value is 0, it means user just ran --watch
			// or explicitly passed 0. We'll default to 5.
			watchInterval = 5
		}

		if watchInterval < 1 {
			ui.DisplayError("Invalid interval", fmt.Errorf("minimum watch interval is 1 minute"))
			os.Exit(1)
		}

		// Setup signal handling for graceful exit
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		ticker := time.NewTicker(time.Duration(watchInterval) * time.Minute)
		defer ticker.Stop()

		// Initial fetch
		ui.DisplayWatchHeader(watchInterval)
		fetchAndDisplayQuota()
		ui.DisplayWatchFooter(time.Now())

		for {
			select {
			case <-sigChan:
				fmt.Println("\nStopping watch mode...")
				return
			case <-ticker.C:
				ui.DisplayWatchHeader(watchInterval)
				fetchAndDisplayQuota()
				ui.DisplayWatchFooter(time.Now())
			}
		}
	}

	fetchAndDisplayQuota()
}

// fetchAndDisplayQuota is the core logic of runQuota separated for watch mode
func fetchAndDisplayQuota() {
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

	// Fetch quota for each account in parallel
	resultsChan := make(chan *ui.AccountQuotaResult, len(accounts))
	var wg sync.WaitGroup

	for _, acc := range accounts {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()

			// Create a new client per goroutine to avoid race conditions
			// as the client holds account-specific state (tokens).
			client := api.NewClient()
			quotaInfo, err := client.GetQuotaInfoForAccount(email)
			if err != nil {
				resultsChan <- &ui.AccountQuotaResult{
					Email: email,
					Error: err.Error(),
				}
				return
			}

			resultsChan <- &ui.AccountQuotaResult{
				Email:        email,
				QuotaSummary: quotaInfo,
			}
		}(acc.Email)
	}

	// Close channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	var quotaResults []*ui.AccountQuotaResult
	for result := range resultsChan {
		quotaResults = append(quotaResults, result)
	}

	// Sort results by email to keep consistent output
	sort.Slice(quotaResults, func(i, j int) bool {
		return quotaResults[i].Email < quotaResults[j].Email
	})

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
		tier := token.TierName
		if tier == "" {
			tier = "Free 📦 (Run 'ag-quota quota' to update)"
		}
		color.Green("✓ Logged in as: %s [%s]", token.Email, tier)
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
	quotaCmd.Flags().IntVarP(&watchInterval, "watch", "w", 0, "Watch quota periodically (default 5m)")
	quotaCmd.Flags().Lookup("watch").NoOptDefVal = "5"
}

func main() {
	// Perform migration if needed (from single-account to multi-account format)
	if err := auth.MigrateIfNeeded(); err != nil {
		fmt.Fprintf(os.Stderr, "Migration warning: %v\n", err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
