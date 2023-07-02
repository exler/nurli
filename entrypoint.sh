#!/bin/sh

set -euo pipefail

cd /app

./nurli migrate

PROCESS_TYPE=$1

echo "Running process type: $PROCESS_TYPE"

if [ "$PROCESS_TYPE" = "serve" ]; then
    ./nurli serve
else
    echo "Unknown process type: $PROCESS_TYPE"
fi
