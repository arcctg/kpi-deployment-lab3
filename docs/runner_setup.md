# Self-Hosted Runner Setup

## Prerequisites

- Ubuntu 24.04 Server VM (runner node, separate from target node)
- Root or sudo access
- GitHub CLI authenticated locally for registration token fetch

## Bootstrap

On the runner VM:

```bash
sudo bash deploy/setup_runner.sh
```

## Register runner manually

Fetch a registration token (do not commit it):

```bash
gh api -X POST repos/arcctg/kpi-deployment-lab3/actions/runners/registration-token
```

On the runner VM:

```bash
cd /opt/actions-runner
./config.sh --url https://github.com/arcctg/kpi-deployment-lab3 --token <REGISTRATION_TOKEN>
sudo ./svc.sh install
sudo ./svc.sh start
```

## SSH access to target node

The bootstrap script generates `/home/runner/.ssh/id_ed25519.pub`. Add it to the deploy user's `~/.ssh/authorized_keys` on the target node.

Store the target connection details in GitHub Secrets:

- `TARGET_HOST`
- `TARGET_USER`
- `TARGET_SSH_KEY`

## After grading

Stop and remove the runner VM to prevent unauthorized use of the self-hosted runner.
