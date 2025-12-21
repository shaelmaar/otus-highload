#!/bin/sh
set -e

CONFIG_DIR=/data
CONFIG_FILE=$CONFIG_DIR/sentinel.conf
BASE_CONFIG=/config/sentinel-base.conf

until getent hosts valkey-master >/dev/null 2>&1; do
  echo "Waiting for valkey-master DNS..."
  sleep 1
done

if [ ! -f "$CONFIG_FILE" ]; then
  cp "$BASE_CONFIG" "$CONFIG_FILE"
fi

exec valkey-sentinel "$CONFIG_FILE"
