#!/bin/sh
# start.sh

host="$1"
shift
port="$1"
shift

until nc -z "$host" "$port"; do
  echo "⏳ Waiting for $host:$port to be ready..."
  sleep 1
done

echo "✅ $host:$port is up. Starting app..."

exec "$@"
