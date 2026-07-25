#!/bin/bash
# Bootstrap script — runs as root on first boot via user_data.
# Installs nginx + aws-cli + certbot, writes a systemd unit that downloads
# the Go binary from S3 (via IAM instance role) on every start, and wires
# nginx as reverse proxy with HTTPS via Let's Encrypt.
#
# Variables injected by Terraform templatefile:
#   ${s3_bucket}, ${s3_object_key}, ${app_name}, ${app_port}, ${domain_name}
# Shell variables use $${...} to escape Terraform template substitution.

set -euo pipefail

# ─── tf vars ───────────────────────────────────────────
S3_BUCKET="${s3_bucket}"
S3_OBJECT_KEY="${s3_object_key}"
APP_NAME="${app_name}"
APP_PORT="${app_port}"
DOMAIN_NAME="${domain_name}"
AWS_REGION="${aws_region}"

# ─── local vars ───────────────────────────────────────────
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
apt-get install -y -qq nginx awscli certbot python3-certbot-nginx

# ─── Create app user + dirs ───────────────────────────────────────────────────
if ! id "$${APP_USER}" &>/dev/null; then
  useradd --system --no-create-home --shell /usr/sbin/nologin "$${APP_USER}"
fi

mkdir -p "$${APP_DIR}"
chown "$${APP_USER}:$${APP_USER}" "$${APP_DIR}"

# ─── systemd unit ─────────────────────────────────────────────────────────────
cat <<EOF > "/etc/systemd/system/$${APP_NAME}.service"
[Unit]
Description=$${APP_NAME} web server
After=network.target

[Service]
Type=simple
User=$${APP_USER}
WorkingDirectory=$${APP_DIR}
Environment=PORT=$${APP_PORT}
ExecStartPre=/usr/bin/aws s3 cp s3://$${S3_BUCKET}/$${S3_OBJECT_KEY} $${APP_BIN} --region $${AWS_REGION}
ExecStartPre=/usr/bin/chmod +x $${APP_BIN}
ExecStart=$${APP_BIN}
Restart=on-failure
RestartSec=30

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable "$${APP_NAME}"

# ─── nginx reverse proxy ──────────────────────────────────────────────────────
if [ -n "$${DOMAIN_NAME}" ]; then
  SERVER_NAME="$${DOMAIN_NAME} www.$${DOMAIN_NAME}"
else
  SERVER_NAME="_"
fi

cat <<EOF > /etc/nginx/sites-available/$${APP_NAME}
server {
    listen 80;
    listen [::]:80;
    server_name $${SERVER_NAME};

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

# ─── HTTPS via certbot ────────────────────────────────────────────────────────
if [ -n "$${DOMAIN_NAME}" ]; then
  certbot --nginx -d "$${DOMAIN_NAME}" -d "www.$${DOMAIN_NAME}" \
    --non-interactive --agree-tos --register-unsafely-without-email \
    --redirect
fi

echo "=== Bootstrap complete ==="
echo "App:  $${APP_BIN} (port $${APP_PORT})"
echo "nginx: :80 -> 127.0.0.1:$${APP_PORT}"
echo "Logs: journalctl -u $${APP_NAME} -f"
