# 📡 Setting Up Telegram Notifications

This guide explains how to configure Telegram notifications for the `ag-quota` CLI to receive real-time updates when your AI model quotas change.

## 🛠️ Step 1: Create a Telegram Bot

1.  Find [@BotFather](https://t.me/botfather) on Telegram.
2.  Send `/newbot` and follow the instructions to get your **Bot Token**.
3.  Save the token securely (e.g., `123456789:ABCDefGhI...`).

## 🆔 Step 2: Get Your Chat ID

1.  Start a conversation with your new bot.
2.  Send any message to the bot.
3.  Use the following API to find your **Chat ID**:
    ```bash
    curl https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates
    ```
    Look for `"chat":{"id":123456789}`. Your ID is `123456789`.

## ⚙️ Step 3: Configure the CLI

Run the following command in your terminal:

```bash
ag-quota config set-telegram --token "YOUR_BOT_TOKEN" --chat-id "YOUR_CHAT_ID"
```

## 🧪 Step 4: Test & Verify

Verify your configuration by sending a test notification with dummy data:

```bash
ag-quota config test-notify
```

You should receive a message in Telegram structured like this:

---

### 🔔 Example Notification Format

**Initial Summary (First Run):**
Sent when you start `--watch` mode to establish a baseline.

> **📊 Quota Summary**
>
> 👤 **user@gmail.com**
>   ✅ **Healthy**
>     - Gemini 3 Flash | 100%
>     - Gemini 3 Pro (Low) | 80%
>   ⚠️ **Warning**
>     - Claude Opus 4.5 | 40%
>   ⛔ **Critical**
>     - Gemini 3 Pro (Thinking) | 10%
>   ❌ **Empty**
>     - GPT-OSS 120B | 0% ⏳ 2h 30m

**Status Update (Subsequent Runs):**
Sent only when a model's status changes (e.g., from Warning to Critical).

> **🔄 Status Update**
>
> 👤 **user@gmail.com**
>   ⛔ **Critical**
>     - Gemini 3 Pro (Thinking) | 10% (↓ 70%)
>   ✅ **Healthy**
>     - GPT-OSS 120B | 100% (↑ 100%)

---

## 💡 Pro Tips

- **Watch Interval**: Use `ag-quota --watch=10` to check every 10 minutes.
- **Multi-Account**: If you use `--all`, notifications will be grouped by account.
- **Silence Modes**: To turn off notifications without removing your token:
  ```bash
  ag-quota config set-telegram --chat-id "" 
  # or simply don't use --watch
  ```

## ⚠️ Safety First

- **Bot Token**: Never share your bot token or commit it to GitHub.
- **Storage**: Your credentials are stored locally in the config directory with restricted permissions (0600).
