#!/usr/bin/env bash
# This script is run as a user_data init script in an EC2 instance

set -euo pipefail
mkdir -p /etc/todoapi
cat << EOF > /etc/todoapi/.env
DB_HOST={db_address}
DB_PORT={db_port}
EOF
systemctl enable --now todoapi
