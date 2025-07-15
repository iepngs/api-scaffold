#!/bin/bash

APP_NAME=${1:-bmstock}
PORT=${2:-9000}
SRC_PATH=/tmp/${APP_NAME}
DEST_PATH=/data/wwwroot/api/${APP_NAME}
CONFIG_PATH=/data/wwwroot/api/config.yaml
LOG_PATH=/dev/null

if [ -f "$SRC_PATH" ]; then
mv "$SRC_PATH" "$DEST_PATH"
fi

chmod +x "$DEST_PATH"

echo "Stopping existing process on port $PORT..."
lsof -ti:$PORT | xargs -r kill || lsof -ti:$PORT | xargs -r kill -9

nohup "$DEST_PATH" -c "$CONFIG_PATH" > "$LOG_PATH" 2>&1 &
sleep 2
lsof -i:$PORT