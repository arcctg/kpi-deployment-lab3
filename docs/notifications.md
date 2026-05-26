# Telegram Notifications

Deploy and verification failures trigger a Telegram message when an annotated tag pipeline fails.

## Setup

1. Open Telegram and message [@BotFather](https://t.me/BotFather).
2. Run `/newbot` and follow the prompts.
3. Copy the bot token and store it as a GitHub Secret:

```bash
gh secret set TELEGRAM_BOT_TOKEN --body "<token>"
```

4. Start a chat with your new bot (send any message).
5. Discover your chat ID:

```bash
curl "https://api.telegram.org/bot<TOKEN>/getUpdates"
```

6. Store the chat ID:

```bash
gh secret set TELEGRAM_CHAT_ID --body "<chat_id>"
```

## Required secrets

- `TELEGRAM_BOT_TOKEN`
- `TELEGRAM_CHAT_ID`

See [SECURITY.md](SECURITY.md) for the full secrets list.
