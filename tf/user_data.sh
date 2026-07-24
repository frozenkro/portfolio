#!/bin/bash
# Bootstrap script — runs as root on first boot via user_data.
# Installs nginx + aws-cli, writes a systemd unit that downloads the Go binary
# from S3 (via IAM instance role) on every start, and wires nginx as reverse proxy.
#
# Variables injected by Terraform templatefile:
#   ${s3_bucket}, ${s3_object_key}, ${app_name}, ${app_port}
# Shell variables use $${...} to escape Terraform template substitution.

set -euo pipefail

# ─── Config (injected by Terraform) ───────────────────────────────────────────
S3_BUCKET="${s3_bucket}"
S3_OBJECT_KEY="${s3_object_key}"
APP_NAME="${app_name}"
APP_PORT="${app_port}"
APP_DIR="/opt/$${APP_NAME}"
APP_BIN="$${APP_DIR}/$${APP_NAME}"
APP_USER="$${APP_NAME}"

# ─── Wait for internet ────────────────────────────────────────────────────────
until curl -sf http://169.254.169.254/latest/meta-data/instance-id >/dev/null 2>&1; do
  sleep 2
done

# ─── Install packages ─────────────────────────────────────────────────────────
export DEBIAN_FRONTEND=noninteractive
apt-get update -qq
apt-get install -y -qq nginx awscli

# ─── Create app user + dirs ───────────────────────────────────────────────────
if ! id "$${APP_USER}" &>/dev/null; then
  useradd --system --no-create-home --shell /usr/sbin/nologin "$${APP_USER}"
fi

mkdir -p "$${APP_DIR}"
chown "$${APP_USER}:$${APP_USER}" "$${APP_DIR}"

# ─── systemd unit ─────────────────────────────────────────────────────────────
# ExecStartPre downloads the binary from S3 on every start attempt.
# If the bucket is empty (e.g. right after terraform apply, before GHA pushes),
# ExecStartPre fails and systemd retries every 30s until the binary appears.
cat <<EOF > "/etc/systemd/system/$${APP_NAME}.service"
[Unit]
Description=$${APP_NAME} web server
After=network.target

[Service]
Type=simple
User=$${APP_USER}
WorkingDirectory=$${APP_DIR}
Environment=PORT=$${APP_PORT}
ExecStartPre=/usr/bin/aws s3 cp s3://$${S3_BUCKET}/$${S3_OBJECT_KEY} $${APP_BIN} --region us-east-1
ExecStart=$${APP_BIN}
Restart=on-failure
RestartSec=30

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable "$${APP_NAME}"

# ─── nginx reverse proxy ──────────────────────────────────────────────────────
cat <<EOF > /etc/nginx/sites-available/$${APP_NAME}
server {
    listen 80;
    listen [::]:80;
    server_name _;

    location / {
        proxy_pass http://127.0.0.1:$${APP_PORT};
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF

ln -sf "/etc/nginx/sites-available/$${APP_NAME}" "/etc/nginx/sites-enabled/$${APP_NAME}"
rm -f /etc/nginx/sites-enabled/default

# ─── Start services ───────────────────────────────────────────────────────────
systemctl restart nginx
systemctl start "$${APP_NAME}"

echo "=== Bootstrap complete ==="
echo "App:  $${APP_BIN} (port $${APP_PORT})"
echo "nginx: :80 -> 127.0.0.1:$${APP_PORT}"
echo "Logs: journalctl -u $${APP_NAME} -f"
