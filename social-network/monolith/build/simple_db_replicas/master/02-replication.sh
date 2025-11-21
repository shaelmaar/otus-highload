#!/bin/bash
set -e

HBA_FILE="$PGDATA/pg_hba.conf"

if ! grep -q "replication.*replicator" "$HBA_FILE"; then
  echo "host replication replicator ${DOCKER_SUBNET} md5" >> "$HBA_FILE"
  echo "Added replication access for replicator on network ${DOCKER_SUBNET}"
fi
