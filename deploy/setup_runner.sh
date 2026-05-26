#!/usr/bin/env bash
set -euo pipefail

RUNNER_DIR="/opt/actions-runner"
RUNNER_VERSION="2.321.0"

apt-get update -qq
DEBIAN_FRONTEND=noninteractive apt-get install -y -qq \
  curl tar openssh-client jq docker.io python3-pip python3-venv gettext-base git

systemctl enable docker
systemctl start docker

mkdir -p "${RUNNER_DIR}"
cd "${RUNNER_DIR}"

if [[ ! -f ./config.sh ]]; then
  curl -sS -o actions-runner.tar.gz -L \
    "https://github.com/actions/runner/releases/download/v${RUNNER_VERSION}/actions-runner-linux-x64-${RUNNER_VERSION}.tar.gz"
  tar xzf actions-runner.tar.gz
  rm actions-runner.tar.gz
  ./bin/installdependencies.sh
fi

useradd -m -s /bin/bash runner 2>/dev/null || true
chown -R runner:runner "${RUNNER_DIR}"

if [[ ! -f /home/runner/.ssh/id_ed25519 ]]; then
  sudo -u runner mkdir -p /home/runner/.ssh
  sudo -u runner ssh-keygen -t ed25519 -N "" -f /home/runner/.ssh/id_ed25519
  chmod 700 /home/runner/.ssh
  chmod 600 /home/runner/.ssh/id_ed25519
  chmod 644 /home/runner/.ssh/id_ed25519.pub
fi

cp /home/runner/.ssh/id_ed25519.pub /vagrant/.runner_pubkey
chmod 644 /vagrant/.runner_pubkey

echo
echo "Runner public key saved to /vagrant/.runner_pubkey"
cat /vagrant/.runner_pubkey
echo

cat <<'EOF'

Runner files are ready in /opt/actions-runner.

Register manually (do NOT commit the token):

  cd /opt/actions-runner
  ./config.sh --url https://github.com/arcctg/kpi-deployment-lab3 --token <REGISTRATION_TOKEN>
  sudo ./svc.sh install runner
  sudo ./svc.sh start

Fetch a fresh registration token locally:

  gh api -X POST repos/arcctg/kpi-deployment-lab3/actions/runners/registration-token

EOF
