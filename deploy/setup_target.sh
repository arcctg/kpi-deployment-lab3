#!/usr/bin/env bash
set -euo pipefail

require_env() {
  local name="$1"
  if [[ -z "${!name:-}" ]]; then
    echo "Missing required environment variable: ${name}" >&2
    exit 1
  fi
}

require_env DB_USER
require_env DB_PASSWORD
require_env DB_NAME
require_env STUDENT_PASSWORD
require_env IMAGE_REF

DEFAULT_PASSWORD="12345678"
REPO_DEPLOY="${REPO_DEPLOY:-/tmp/mywebapp-deploy}"

user_exists() {
  id "$1" &>/dev/null
}

apt-get update -qq
DEBIAN_FRONTEND=noninteractive apt-get install -y -qq \
  docker.io nginx postgresql curl ca-certificates gettext-base

systemctl enable docker
systemctl start docker

if ! user_exists student; then
  useradd -m -s /bin/bash student
fi
usermod -aG sudo student
echo "student:${STUDENT_PASSWORD}" | chpasswd

if ! user_exists teacher; then
  useradd -m -s /bin/bash teacher
fi
usermod -aG sudo teacher
echo "teacher:${DEFAULT_PASSWORD}" | chpasswd
chage -d 0 teacher

if ! user_exists app; then
  useradd -r -s /usr/sbin/nologin app
fi

if ! getent group operator >/dev/null; then
  groupadd operator
fi
if ! user_exists operator; then
  useradd -m -s /bin/bash -g operator operator
fi
echo "operator:${DEFAULT_PASSWORD}" | chpasswd
chage -d 0 operator

systemctl enable postgresql
systemctl start postgresql

if [[ "$(sudo -u postgres psql -tAc "SELECT 1 FROM pg_roles WHERE rolname='${DB_USER}'")" != "1" ]]; then
  sudo -u postgres psql -c "CREATE ROLE ${DB_USER} WITH LOGIN PASSWORD '${DB_PASSWORD}'"
fi
if [[ "$(sudo -u postgres psql -tAc "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}'")" != "1" ]]; then
  sudo -u postgres createdb -O "${DB_USER}" "${DB_NAME}"
fi

mkdir -p /etc/mywebapp
envsubst < "${REPO_DEPLOY}/config.yaml.tmpl" > /etc/mywebapp/config.yaml
chown root:app /etc/mywebapp/config.yaml
chmod 640 /etc/mywebapp/config.yaml

export IMAGE_REF
APP_UID=$(id -u app)
APP_GID=$(id -g app)
export APP_UID APP_GID IMAGE_REF
envsubst < "${REPO_DEPLOY}/mywebapp.service.tmpl" > /etc/systemd/system/mywebapp.service

rm -f /etc/nginx/sites-enabled/default
install -m 644 "${REPO_DEPLOY}/nginx.conf" /etc/nginx/sites-available/mywebapp
ln -sf /etc/nginx/sites-available/mywebapp /etc/nginx/sites-enabled/mywebapp
nginx -t

install -m 440 "${REPO_DEPLOY}/sudoers-operator" /etc/sudoers.d/operator
chown root:root /etc/sudoers.d/operator

install -m 440 "${REPO_DEPLOY}/sudoers-student" /etc/sudoers.d/student
chown root:root /etc/sudoers.d/student

mkdir -p /home/student
echo "9" > /home/student/gradebook
chown student:student /home/student/gradebook

systemctl daemon-reload
systemctl enable nginx mywebapp.service
systemctl restart nginx
systemctl restart mywebapp.service

echo "Target node setup complete"
