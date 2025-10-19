#!/bin/bash

until /usr/bin/mongosh --port 27017 --quiet --eval 'db.getMongo()';
do
  sleep 1
done

readarray -d';' -t shards <<< "$SHARDS;"

for shard in "${shards[@]}"
do
  if [ $(echo "$shard" | wc -l) -gt 1 ];
  then
    break
  fi
  /usr/bin/mongosh --eval "sh.addShard(\"$shard\")"
done

APP_DB=${MONGO_DB_NAME:-appdb}
APP_USER=${MONGO_DB_USER:-appuser}
APP_PASS=${MONGO_DB_PASSWORD:-apppass}

echo "🗄️  Initializing app database $APP_DB..."
/usr/bin/mongosh --quiet --port 27017 <<EOF
use $APP_DB;

print("✅ Database '$APP_DB' created successfully");
EOF

echo "🚀 Router initialization complete"