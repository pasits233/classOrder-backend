#!/bin/bash

PROJECT_DIR=$(cd "$(dirname "$0")"; pwd)
BIN="$PROJECT_DIR/bin/backend"
LOG="$PROJECT_DIR/backend.log"
PID_FILE="$PROJECT_DIR/backend.pid"

start() {
    if [ -f "$PID_FILE" ] && kill -0 $(cat "$PID_FILE") 2>/dev/null; then
        echo "Backend is already running (PID: $(cat $PID_FILE))"
        exit 0
    fi
    echo "Starting backend..."
    nohup "$BIN" > "$LOG" 2>&1 &
    echo $! > "$PID_FILE"
    echo "Backend started (PID: $(cat $PID_FILE))"
}

stop() {
    if [ -f "$PID_FILE" ] && kill -0 $(cat "$PID_FILE") 2>/dev/null; then
        echo "Stopping backend (PID: $(cat $PID_FILE))..."
        kill $(cat "$PID_FILE")
        rm -f "$PID_FILE"
        echo "Stopped."
    else
        echo "Backend is not running."
    fi
}

restart() {
    stop
    sleep 1
    start
}

status() {
    if [ -f "$PID_FILE" ] && kill -0 $(cat "$PID_FILE") 2>/dev/null; then
        echo "Backend is running (PID: $(cat $PID_FILE))"
    else
        echo "Backend is not running."
    fi
}

log() {
    tail -f "$LOG"
}

case "$1" in
    start) start ;;
    stop) stop ;;
    restart) restart ;;
    status) status ;;
    log) log ;;
    *)
        echo "Usage: $0 {start|stop|restart|status|log}"
        exit 1
        ;;
esac 