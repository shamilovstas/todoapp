#!/usr/bin/env bash

set -euo pipefail
source /etc/todoapi/.env
migrate -database postgres://$DB_USER:$DB_PASS@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable -path /usr/share/todoapi/sql up
