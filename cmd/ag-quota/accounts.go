package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gundamkid/anti-gravity-quota/internal/auth"
	"github.com/gundamkid/anti-gravity-quota/internal/ui"
	"github.com/spf13/cobra"
)

// accountsCmd represents the accounts command
var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage saved accounts",
	Long: `List, switch between, and manage your saved Google accounts.
	
Examples:
  ag-quota accounts list              # List all saved accounts
  ag-quota accounts switch user@gmail.com   # Switch to another account
  ag-quota accounts set-default user@gmail.com  # Set default account`,
	RunE: runAccountsList, // Default action: list accounts
}

// accountsListCmd represents the accounts list command
var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved accounts",
	Long:  `Display all saved Google accounts with their status and default marker.`,
	RunE:  runAccountsList,
}

// accountsSwitchCmd represents the accounts switch command
var accountsSwitchCmd = &cobra.Command{
	Use:   "switch <email>",
	Short: "Switch to another account",
	Long:  `Switch to another saved account and set it as the default.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runAccountsSwitch,
}

// accountsSetDefaultCmd represents the accounts set-default command (alias for switch)
var accountsSetDefaultCmd = &cobra.Command{
	Use:   "set-default <email>",
	Short: "Set the default account",
	Long:  `Set a saved account as the default account (alias for switch).`,
	Args:  cobra.ExactArgs(1),
	RunE:  runAccountsSwitch, // Same implementation as switch
}

// runAccountsList handles listing all saved accounts
func runAccountsList(cmd *cobra.Command, args []string) error {
	mgr, err := auth.NewAccountManager()
	if err != nil {
		return fmt.Errorf("failed to initialize account manager: %w", err)
	}

	accounts, err := mgr.ListAccounts()
	if err != nil {
		return fmt.Errorf("failed to list accounts: %w", err)
	}

	if len(accounts) == 0 {
		fmt.Println()
		color.Yellow("No accounts found")
		fmt.Println()
		fmt.Println("Run 'ag-quota login' to add your first account.")
		fmt.Println()
		return nil
	}

	ui.DisplayAccountsList(accounts)
	return nil
}

// runAccountsSwitch handles switching to another account
func runAccountsSwitch(cmd *cobra.Command, args []string) error {
	email := args[0]

	mgr, err := auth.NewAccountManager()
	if err != nil {
		return fmt.Errorf("failed to initialize account manager: %w", err)
	}

	// Set the account as default
	if err := mgr.SetDefaultAccount(email); err != nil {
		// Check if it's a "not found" error
		if os.IsNotExist(err) {
			fmt.Println()
			color.Red("✗ Account not found: %s", email)
			fmt.Println()
			fmt.Println("Available accounts:")

			// Show available accounts
			accounts, listErr := mgr.ListAccounts()
			if listErr == nil && len(accounts) > 0 {
				for _, acc := range accounts {
					fmt.Printf("  • %s\n", acc.Email)
				}
			} else {
				fmt.Println("  (none)")
			}

			fmt.Println()
			fmt.Println("Run 'ag-quota login' to add a new account.")
			fmt.Println()
			return nil
		}
		return fmt.Errorf("failed to switch account: %w", err)
	}

	fmt.Println()
	color.Green("✓ Switched to %s", email)
	fmt.Println()

	return nil
}

func init() {
	// Add accounts command to root
	rootCmd.AddCommand(accountsCmd)

	// Add subcommands to accounts
	accountsCmd.AddCommand(accountsListCmd)
	accountsCmd.AddCommand(accountsSwitchCmd)
	accountsCmd.AddCommand(accountsSetDefaultCmd)
}
