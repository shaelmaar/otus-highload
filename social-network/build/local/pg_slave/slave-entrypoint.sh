#!/bin/bash
set -e

if [ ! -f /var/lib/postgresql/.standby_initialized ]; then
  echo "Initializing standby from master..."
  rm -rf /var/lib/postgresql/data/*

  PGPASSWORD=repl_secret pg_basebackup -h monolith_db_master -D /var/lib/postgresql/data -U replicator -Fp -Xs -P -R

  touch /var/lib/postgresql/.standby_initialized
fi

exec docker-entrypoint.sh "$@"
