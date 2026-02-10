package main

import (
	"context"
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
	"github.com/gundamkid/anti-gravity-quota/internal/notify"
	"github.com/gundamkid/anti-gravity-quota/internal/ui"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	version       = "0.1.1"
	jsonOutput    bool
	accountFlag   string
	allFlag       bool
	watchInterval int
	compactFlag   bool
	noCompactFlag bool

	// Notifications
	notifRegistry *notify.Registry
	stateTracker  *notify.StateTracker
	msgFormatter  *notify.MessageFormatter
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
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		runQuota(ctx, cmd, args)
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
func runQuota(ctx context.Context, cmd *cobra.Command, args []string) {
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
		fetchAndDisplayQuota(ctx, true)
		ui.DisplayWatchFooter(time.Now())

		for {
			select {
			case <-sigChan:
				fmt.Println("\nStopping watch mode...")
				return
			case <-ticker.C:
				ui.DisplayWatchHeader(watchInterval)
				fetchAndDisplayQuota(ctx, true)
				ui.DisplayWatchFooter(time.Now())
			}
		}
	}

	fetchAndDisplayQuota(ctx, false)
}

// fetchAndDisplayQuota is the core logic of runQuota separated for watch mode
func fetchAndDisplayQuota(ctx context.Context, triggerNotify bool) {
	var finalResults []*ui.AccountQuotaResult

	// Handle --all flag
	if allFlag {
		finalResults = runQuotaForAllAccounts(ctx)
	} else {
		// Handle --account flag or default
		email := accountFlag
		if email == "" {
			token, err := auth.LoadToken()
			if err != nil {
				if jsonOutput {
					fmt.Fprintf(os.Stderr, `{"error": "not logged in"}%s`, "\n")
				} else {
					ui.DisplayNotLoggedIn()
				}
				os.Exit(1)
			}
			email = token.Email
		}

		res, err := fetchQuotaForAccountResult(ctx, email)
		if err != nil {
			if ctx.Err() == nil {
				if jsonOutput {
					fmt.Fprintf(os.Stderr, `{"error": "failed to fetch quota", "message": "%s"}%s`, err.Error(), "\n")
				} else {
					ui.DisplayError("Failed to fetch quota information", err)
				}
				os.Exit(1)
			}
			return
		}
		finalResults = []*ui.AccountQuotaResult{res}
	}

	if finalResults == nil {
		return
	}

	// Determine if compact mode should be used
	displayOpts := ui.DisplayOptions{
		Compact: false,
	}

	if compactFlag {
		displayOpts.Compact = true
	} else if !noCompactFlag {
		// Auto-detect based on terminal width
		width := ui.GetTerminalWidth()
		if width < 80 {
			displayOpts.Compact = true
		}
	}

	// Display results
	if jsonOutput {
		if allFlag {
			if err := ui.DisplayAllAccountsQuotaJSON(finalResults); err != nil {
				fmt.Fprintf(os.Stderr, "Error displaying JSON: %v\n", err)
				os.Exit(1)
			}
		} else if len(finalResults) > 0 {
			if err := ui.DisplayQuotaSummaryJSON(finalResults[0].QuotaSummary); err != nil {
				fmt.Fprintf(os.Stderr, "Error displaying JSON: %v\n", err)
				os.Exit(1)
			}
		}
	} else {
		if allFlag {
			ui.DisplayAllAccountsQuota(finalResults, displayOpts)
		} else if len(finalResults) > 0 {
			ui.DisplayQuotaSummary(finalResults[0].QuotaSummary, displayOpts)
		}
	}

	// Handle notifications if enabled and requested
	if triggerNotify && notifRegistry != nil {
		var allChanges []notify.StatusChange
		for _, res := range finalResults {
			if res.QuotaSummary != nil {
				changes := stateTracker.Update(res.Email, res.QuotaSummary.Models)
				allChanges = append(allChanges, changes...)
			}
		}

		if len(allChanges) > 0 {
			msg := msgFormatter.FormatChanges(allChanges)
			notifRegistry.NotifyAll(ctx, msg)
		}
	}
}

// fetchQuotaForAccountResult is a helper to fetch quota and return as result struct
func fetchQuotaForAccountResult(ctx context.Context, email string) (*ui.AccountQuotaResult, error) {
	if !jsonOutput && !allFlag {
		fmt.Printf("Fetching quota for %s... ", email)
	}

	client := api.NewClient()
	quotaInfo, err := client.GetQuotaInfoForAccount(ctx, email)
	if err != nil {
		return nil, err
	}

	if !jsonOutput && !allFlag {
		fmt.Println(color.GreenString("✓"))
	}

	return &ui.AccountQuotaResult{
		Email:        email,
		QuotaSummary: quotaInfo,
	}, nil
}

// runQuotaForAllAccounts fetches and returns quota for all saved accounts
func runQuotaForAllAccounts(ctx context.Context) []*ui.AccountQuotaResult {
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

	if !jsonOutput && !allFlag { // Redundant but for clarity
		fmt.Println()
		fmt.Printf("Fetching quota for %d account(s)...\n", len(accounts))
		fmt.Println()
	}

	// Fetch quota for each account in parallel using errgroup
	quotaResults := make([]*ui.AccountQuotaResult, len(accounts))
	g, gCtx := errgroup.WithContext(ctx)
	var mu sync.Mutex

	for i, acc := range accounts {
		idx, email := i, acc.Email
		g.Go(func() error {
			// Create a new client per goroutine to avoid race conditions
			client := api.NewClient()
			quotaInfo, apiErr := client.GetQuotaInfoForAccount(gCtx, email)

			mu.Lock()
			defer mu.Unlock()

			if apiErr != nil {
				// Record individual account error instead of returning it to errgroup
				// This prevents one bad account from stopping the entire --all fetch.
				quotaResults[idx] = &ui.AccountQuotaResult{
					Email: email,
					Error: apiErr.Error(),
				}
				return nil
			}

			quotaResults[idx] = &ui.AccountQuotaResult{
				Email:        email,
				QuotaSummary: quotaInfo,
			}
			return nil
		})
	}

	// Wait for completion. Fatal errors (cancellation) will still cause Wait to return error
	err = g.Wait()
	if ctx.Err() != nil {
		return nil
	}

	if err != nil {
		// For other fatal errors that might still propagate
		if jsonOutput {
			fmt.Fprintf(os.Stderr, `{"error": "fatal error", "message": "%s"}%s`, err.Error(), "\n")
		} else {
			ui.DisplayError("fatal error during fetch", err)
		}
		os.Exit(1)
	}

	// Filter out nil results (though with g.Wait they should all be filled if no error)
	var finalResults []*ui.AccountQuotaResult
	for _, r := range quotaResults {
		if r != nil {
			finalResults = append(finalResults, r)
		}
	}

	sort.Slice(finalResults, func(i, j int) bool {
		return finalResults[i].Email < finalResults[j].Email
	})

	return finalResults
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
	rootCmd.AddCommand(configCmd)

	// Add flags to root asPersistentFlags so they are available to subcommands and when running root directly
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")
	rootCmd.PersistentFlags().StringVar(&accountFlag, "account", "", "Check quota for specific account")
	rootCmd.PersistentFlags().BoolVar(&allFlag, "all", false, "Check quota for all accounts")
	rootCmd.PersistentFlags().IntVarP(&watchInterval, "watch", "w", 0, "Watch quota periodically (default 5m)")
	rootCmd.PersistentFlags().BoolVar(&compactFlag, "compact", false, "Force compact mode display")
	rootCmd.PersistentFlags().BoolVar(&noCompactFlag, "no-compact", false, "Force full mode display (disable auto-compact)")
	rootCmd.PersistentFlags().Lookup("watch").NoOptDefVal = "5"
}

func main() {
	// Initialize notifications
	initNotifications()

	// Perform migration if needed (from single-account to multi-account format)
	if err := auth.MigrateIfNeeded(); err != nil {
		fmt.Fprintf(os.Stderr, "Migration warning: %v\n", err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func initNotifications() {
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	if !cfg.Notifications.Enabled {
		return
	}

	notifRegistry = notify.NewRegistry()
	stateTracker = notify.NewStateTracker()
	msgFormatter = notify.NewMessageFormatter()

	// Register Telegram if configured
	if cfg.Notifications.Telegram.BotToken != "" && cfg.Notifications.Telegram.ChatID != "" {
		notifRegistry.Register(notify.NewTelegramNotifier(
			cfg.Notifications.Telegram.BotToken,
			cfg.Notifications.Telegram.ChatID,
		))
	}
}
