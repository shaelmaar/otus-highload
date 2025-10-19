#!/bin/bash
set -e

echo "Waiting mongos..."

until mongosh --host mongos1:27017 --eval "db.runCommand({ ping: 1 })" --quiet; do
  sleep 2
done

echo "Mongos is available."

mongosh --host mongos1:27017 <<EOF
  sh.enableSharding("dialogs_db");
  sh.shardCollection("dialogs_db.dialogMessages", {"dialogID": "hashed"});
EOF

echo "Sharding enabled"