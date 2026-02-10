# Setting Up Telegram Bot for AG-Quota Notifications

This guide will walk you through setting up a Telegram bot to receive quota status notifications from the `ag-quota` CLI when using watch mode.

## Prerequisites

- You need a Telegram account
- The `ag-quota` CLI installed and configured (run `ag-quota login` first)
- At least one Google account authenticated

---

## Step 1: Create a Telegram Bot

1. **Open Telegram** and search for `@BotFather` (the official Telegram bot for managing bots)

2. **Start a chat** with BotFather by clicking "START" or typing `/start`

3. **Create a new bot** by sending the command:
   ```
   /newbot
   ```

4. **Choose a name** for your bot when BotFather asks. For example:
   ```
   AG Quota Notifier
   ```

5. **Choose a username** for your bot (must end in `bot`). For example:
   ```
   ag_quota_notifier_bot
   ```

6. **Save the Bot Token**. BotFather will reply with a message containing your bot token. It looks like this:
   ```
   123456789:ABCdefGHIjklMNOpqrsTUVwxyz-1234567890
   ```
   
   ⚠️ **IMPORTANT**: Keep this token secure! Anyone with this token can control your bot.

---

## Step 2: Get Your Chat ID

1. **Start a chat** with your newly created bot:
   - Find your bot in Telegram by searching for its username (e.g., `@ag_quota_notifier_bot`)
   - Click "START" to initiate the conversation

2. **Send any message** to the bot (e.g., "Hello")

3. **Get your Chat ID** using one of these methods:

   **Method A: Using a web browser**
   - Open this URL in your browser, replacing `YOUR_BOT_TOKEN` with your actual token:
     ```
     https://api.telegram.org/botYOUR_BOT_TOKEN/getUpdates
     ```
   - Look for the `"chat":{"id":` field in the JSON response
   - Your Chat ID will be a number like `123456789`

   **Method B: Using curl command**
   ```bash
   curl https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates
   ```
   - Look for `"chat":{"id":` in the output

4. **Save your Chat ID** - you'll need it in the next step

---

## Step 3: Configure ag-quota with Telegram Credentials

Now that you have both the Bot Token and Chat ID, configure `ag-quota`:

```bash
ag-quota config set-telegram --token "YOUR_BOT_TOKEN" --chat-id "YOUR_CHAT_ID"
```

**Example:**
```bash
ag-quota config set-telegram --token "123456789:ABCdefGHIjklMNOpqrsTUVwxyz-1234567890" --chat-id "987654321"
```

The CLI will:
- ✅ Validate your bot token by calling the Telegram API
- ✅ Save your credentials securely (with 0600 permissions)
- ✅ Automatically enable notifications

---

## Step 4: Verify Configuration

Check that your Telegram settings are configured correctly:

```bash
ag-quota config get-telegram
```

You should see output like:
```
Telegram Configuration:
  Bot Token: 123456***890 (masked)
  Chat ID: 987654321
  Status: ✓ Configured
```

---

## Step 5: Test Notifications with Watch Mode

Now test the notification system by running watch mode:

```bash
ag-quota quota --watch 1
```

This will:
- Check your quota every 1 minute
- Detect any status changes (HEALTHY → WARNING → CRITICAL → EMPTY)
- Send Telegram notifications when changes occur

### What to Expect:

**On first run**, you'll receive a notification for any models that are NOT in HEALTHY status. For example:
```
🔔 [AG-Quota] Status Update

🔴 *CRITICAL*
• Gemini 3 Flash: 5% - Reset in 2h 30m

⚠️ *WARNING*
• Claude 4 Opus: 45%
```

**On subsequent runs**, you'll only receive notifications when status changes:
```
🔔 [AG-Quota] Status Update

🚫 *EMPTY*
• Gemini 3 Flash: 0% (5% ↓ 0% (↓5%)) - Reset in 1h 45m
```

---

## Advanced: Testing Multiple Accounts

If you have multiple Google accounts configured, you can watch all of them:

```bash
ag-quota quota --all --watch 5
```

Notifications will include the account email:
```
🔔 [AG-Quota] Status Update

🔴 *CRITICAL*
• Claude 4 Sonnet: 12% (80% ↓ 12% (↓68%))
  └ Account: work@gmail.com
```

---

## Troubleshooting

### Issue: "Failed to validate bot token"
- **Solution**: Double-check your bot token. Make sure there are no extra spaces or characters.

### Issue: "No notifications received"
- **Check 1**: Make sure you've started a chat with your bot (sent at least one message)
- **Check 2**: Verify notifications are enabled: `ag-quota config get-telegram`
- **Check 3**: Test with a short interval: `ag-quota quota --watch 1`

### Issue: "Rate limit exceeded"
- **Explanation**: The CLI limits notifications to 10 messages per minute to prevent API abuse
- **Solution**: Wait 1 minute before the next notification batch

### Issue: "Unauthorized" error from Telegram API
- **Solution**: Your bot token may be invalid. Create a new bot and update the configuration.

---

## Security Notes

⚠️ **Important Security Considerations:**

1. **Never share your bot token** - Treat it like a password
2. **Credentials are stored** in `~/.config/ag-quota/config.json` with 0600 permissions
3. **Only you** can send messages to your bot by default
4. If you suspect your token was compromised:
   - Talk to `@BotFather` in Telegram
   - Send `/revoke` to revoke the old token
   - Create a new bot and reconfigure ag-quota

---

## Disabling Notifications

To temporarily disable notifications without removing your credentials:

```bash
ag-quota config set-notifications --enabled=false
```

To re-enable:

```bash
ag-quota config set-notifications --enabled=true
```

---

## Next Steps

- **Customize watch interval**: Adjust the `--watch` flag value (default: 5 minutes)
- **Run in background**: Use `nohup` or `screen` to keep watch mode running
- **Monitor specific accounts**: Use `--account` flag to watch a single account

For more information, run:
```bash
ag-quota --help
ag-quota config --help
```

---

**Happy monitoring! 🚀**
