package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gundamkid/anti-gravity-quota/internal/config"
	"github.com/gundamkid/anti-gravity-quota/internal/notify"
	"github.com/gundamkid/anti-gravity-quota/internal/ui"
	"github.com/spf13/cobra"
)

var (
	telegramToken  string
	telegramChatID string
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage application configuration",
	Long:  `View and modify application settings like notification credentials.`,
}

// setTelegramCmd represents the set-telegram command
var setTelegramCmd = &cobra.Command{
	Use:   "set-telegram",
	Short: "Set Telegram notification credentials",
	Long:  `Configure the Telegram bot token and chat ID for notifications.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			ui.DisplayError("Failed to load config", err)
			os.Exit(1)
		}

		updated := false
		if telegramToken != "" {
			// Validate token before saving
			tn := cfg.Notifications.Telegram
			tn.BotToken = telegramToken

			tNotifier := notify.NewTelegramNotifier(tn.BotToken, "")
			fmt.Print("Validating Telegram Bot Token... ")
			if err := tNotifier.Validate(cmd.Context()); err != nil {
				fmt.Println(color.RedString("FAILED"))
				ui.DisplayError("Invalid Telegram token", err)
				os.Exit(1)
			}
			fmt.Println(color.GreenString("OK"))

			cfg.Notifications.Telegram.BotToken = telegramToken
			updated = true
		}
		if telegramChatID != "" {
			cfg.Notifications.Telegram.ChatID = telegramChatID
			updated = true
		}

		if !updated {
			color.Yellow("No changes provided. Use --token and --chat-id flags.")
			return
		}

		// Enable notifications if credentials are set
		if cfg.Notifications.Telegram.BotToken != "" && cfg.Notifications.Telegram.ChatID != "" {
			cfg.Notifications.Enabled = true
		}

		if err := config.SaveConfig(cfg); err != nil {
			ui.DisplayError("Failed to save config", err)
			os.Exit(1)
		}

		color.Green("✓ Telegram configuration updated successfully")
		if cfg.Notifications.Enabled {
			color.Cyan("Notifications are now ENABLED")
		}
	},
}

// getTelegramCmd represents the get-telegram command
var getTelegramCmd = &cobra.Command{
	Use:   "get-telegram",
	Short: "View Telegram notification credentials",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			ui.DisplayError("Failed to load config", err)
			os.Exit(1)
		}

		fmt.Println("Telegram Configuration")
		fmt.Println("======================")
		fmt.Printf("Notifications: %s\n", statusString(cfg.Notifications.Enabled))
		fmt.Printf("Bot Token:     %s\n", maskToken(cfg.Notifications.Telegram.BotToken))
		fmt.Printf("Chat ID:       %s\n", cfg.Notifications.Telegram.ChatID)
	},
}

// testNotifyCmd represents the test-notify command
var testNotifyCmd = &cobra.Command{
	Use:   "test-notify",
	Short: "Send a test notification to verified providers",
	Run: func(cmd *cobra.Command, args []string) {
		if notifRegistry == nil || len(notifRegistry.List()) == 0 {
			color.Yellow("No notification providers are registered or enabled.")
			fmt.Println("Use 'ag-quota config set-telegram' to configure Telegram.")
			return
		}

		fmt.Println("Sending test notification...")
		msg := notify.Message{
			Title:    "Test Notification 🚀",
			Body:     "This is a test notification from Anti-Gravity Quota CLI. Your configuration is working correctly!",
			Severity: notify.SeverityInfo,
		}

		errs := notifRegistry.NotifyAll(cmd.Context(), msg)
		if len(errs) > 0 {
			for _, err := range errs {
				ui.DisplayError("Failed to send notification", err)
			}
			os.Exit(1)
		}

		color.Green("✓ Test notification sent successfully!")
	},
}

func statusString(enabled bool) string {
	if enabled {
		return color.GreenString("ENABLED")
	}
	return color.RedString("DISABLED")
}

func maskToken(token string) string {
	if token == "" {
		return color.HiBlackString("not set")
	}
	if len(token) <= 10 {
		return "**********"
	}
	return token[:6] + "..." + token[len(token)-4:]
}

func init() {
	// Add subcommands to config
	configCmd.AddCommand(setTelegramCmd)
	configCmd.AddCommand(getTelegramCmd)
	configCmd.AddCommand(testNotifyCmd)

	// Add flags to set-telegram
	setTelegramCmd.Flags().StringVar(&telegramToken, "token", "", "Telegram bot token")
	setTelegramCmd.Flags().StringVar(&telegramChatID, "chat-id", "", "Telegram chat ID")
}
