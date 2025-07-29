#!/bin/sh
# start.sh

set -e

docker compose pull
docker compose build
docker compose up "$@"