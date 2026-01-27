package main

import (
	"errors"
	"fmt"

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
  ag-quota accounts default user@gmail.com  # Set default account
  ag-quota accounts switch user@gmail.com   # Alias for default
  ag-quota accounts remove user@gmail.com   # Remove account`,
	RunE: runAccountsList, // Default action: list accounts
}

// accountsListCmd represents the accounts list command
var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved accounts",
	Long:  `Display all saved Google accounts with their status and default marker.`,
	RunE:  runAccountsList,
}

// accountsDefaultCmd represents the accounts default command
var accountsDefaultCmd = &cobra.Command{
	Use:     "default <email>",
	Aliases: []string{"switch"},
	Short:   "Set the default account",
	Long:    `Set a saved account as the default account. Commands will use this account if none is specified.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runAccountsDefault,
}

// accountsRemoveCmd represents the accounts remove command
var accountsRemoveCmd = &cobra.Command{
	Use:   "remove <email>",
	Short: "Remove a saved account",
	Long:  `Delete the saved authentication token for the specified account.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runAccountsRemove,
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

// runAccountsDefault handles setting the default account
func runAccountsDefault(cmd *cobra.Command, args []string) error {
	email := args[0]

	mgr, err := auth.NewAccountManager()
	if err != nil {
		return fmt.Errorf("failed to initialize account manager: %w", err)
	}

	// Set the account as default
	if err := mgr.SetDefaultAccount(email); err != nil {
		// Check if it's a "not found" error
		if errors.Is(err, auth.ErrAccountNotFound) {
			fmt.Printf("Account %s not found.\n", email)
			return nil
		}
		return fmt.Errorf("failed to set default account: %w", err)
	}

	fmt.Printf("✅ Default account set to: %s\n", email)
	return nil
}

// runAccountsRemove handles removing a saved account
func runAccountsRemove(cmd *cobra.Command, args []string) error {
	email := args[0]

	mgr, err := auth.NewAccountManager()
	if err != nil {
		return fmt.Errorf("failed to initialize account manager: %w", err)
	}

	// For simplicity, we skip confirmation for now, or we could add a --force flag.
	// But let's follow standard CLI practice if we can.
	// Simple confirmation for now.
	fmt.Printf("Are you sure you want to remove account %s? (y/N): ", email)
	var response string
	fmt.Scanln(&response)
	if response != "y" && response != "Y" {
		fmt.Println("Aborted.")
		return nil
	}

	if err := mgr.RemoveAccount(email); err != nil {
		if errors.Is(err, auth.ErrAccountNotFound) {
			fmt.Printf("Account %s not found.\n", email)
			return nil
		}
		return fmt.Errorf("failed to remove account: %w", err)
	}

	fmt.Printf("✅ Account %s removed successfully.\n", email)
	return nil
}

func init() {
	// Add accounts command to root
	rootCmd.AddCommand(accountsCmd)

	// Add subcommands to accounts
	accountsCmd.AddCommand(accountsListCmd)
	accountsCmd.AddCommand(accountsDefaultCmd)
	accountsCmd.AddCommand(accountsRemoveCmd)
}
