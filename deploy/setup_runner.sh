#!/usr/bin/env bash
set -euo pipefail

RUNNER_DIR="/opt/actions-runner"
RUNNER_VERSION="2.321.0"

apt-get update -qq
DEBIAN_FRONTEND=noninteractive apt-get install -y -qq \
  curl tar openssh-client jq docker.io python3-pip python3-venv gettext-base

systemctl enable docker
systemctl start docker

mkdir -p "${RUNNER_DIR}"
cd "${RUNNER_DIR}"

if [[ ! -f ./config.sh ]]; then
  curl -sS -o actions-runner.tar.gz -L \
    "https://github.com/actions/runner/releases/download/v${RUNNER_VERSION}/actions-runner-linux-x64-${RUNNER_VERSION}.tar.gz"
  tar xzf actions-runner.tar.gz
  rm actions-runner.tar.gz
fi

if [[ ! -f /home/runner/.ssh/id_ed25519 ]]; then
  useradd -m -s /bin/bash runner 2>/dev/null || true
  sudo -u runner mkdir -p /home/runner/.ssh
  sudo -u runner ssh-keygen -t ed25519 -N "" -f /home/runner/.ssh/id_ed25519
  chmod 700 /home/runner/.ssh
  chmod 600 /home/runner/.ssh/id_ed25519
  chmod 644 /home/runner/.ssh/id_ed25519.pub
  echo
  echo "Add this public key to the target node authorized_keys for the deploy user:"
  cat /home/runner/.ssh/id_ed25519.pub
  echo
fi

cat <<'EOF'

Runner files are ready in /opt/actions-runner.

Register manually (do NOT commit the token):

  cd /opt/actions-runner
  ./config.sh --url https://github.com/arcctg/kpi-deployment-lab3 --token <REGISTRATION_TOKEN>
  sudo ./svc.sh install
  sudo ./svc.sh start

Fetch a fresh registration token locally:

  gh api -X POST repos/arcctg/kpi-deployment-lab3/actions/runners/registration-token

EOF
