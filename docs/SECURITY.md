# Security and Secrets

This repository must not contain production credentials. Use GitHub Secrets for deployment and a local `.env` file (gitignored) for development.

## GitHub Secrets

| Secret | Purpose |
|---|---|
| `TARGET_HOST` | Target node IP or hostname reachable from the self-hosted runner |
| `TARGET_USER` | SSH user on the target node used by deploy jobs |
| `TARGET_SSH_KEY` | Private SSH key used by the runner to reach the target node |
| `DB_PASSWORD` | PostgreSQL role password for production |
| `STUDENT_PASSWORD` | Initial password for the `student` user on the target node |
| `TELEGRAM_BOT_TOKEN` | Telegram bot token for failure notifications |
| `TELEGRAM_CHAT_ID` | Telegram chat ID for failure notifications |

`GITHUB_TOKEN` is provided automatically by GitHub Actions for GHCR push.

## Local development

Copy `.env.example` to `.env` and set local values:

```bash
cp .env.example .env
```

Render config templates locally:

```bash
export $(grep -v '^#' .env | xargs)
envsubst < deploy/docker-config.yaml.tmpl > deploy/rendered/docker-config.yaml
```

## Runner registration

Self-hosted runner registration tokens must never be committed. Fetch a fresh token when needed:

```bash
gh api -X POST repos/OWNER/REPO/actions/runners/registration-token
```

Run `./config.sh` on the runner VM manually with that token.
