#!/bin/bash
# Script to stop the ONVIF backend and frontend services
# Author: GitHub Copilot
# Created: May 6, 2025

BASE_DIR=$(dirname "$(readlink -f "$0")")
BACKEND_PID_FILE="$BASE_DIR/logs/backend.pid"
FRONTEND_PID_FILE="$BASE_DIR/logs/frontend.pid"

echo "Stopping ONVIF services..."

# Stop backend if running
if [ -f "$BACKEND_PID_FILE" ]; then
    BACKEND_PID=$(cat "$BACKEND_PID_FILE")
    if ps -p $BACKEND_PID > /dev/null; then
        echo "Stopping backend (PID: $BACKEND_PID)..."
        kill $BACKEND_PID
        echo "Backend stopped."
    else
        echo "Backend process is not running."
    fi
    rm -f "$BACKEND_PID_FILE"
else
    echo "Backend PID file not found."
fi

# Stop frontend if running
if [ -f "$FRONTEND_PID_FILE" ]; then
    FRONTEND_PID=$(cat "$FRONTEND_PID_FILE")
    if ps -p $FRONTEND_PID > /dev/null; then
        echo "Stopping frontend (PID: $FRONTEND_PID)..."
        kill $FRONTEND_PID
        echo "Frontend stopped."
    else
        echo "Frontend process is not running."
    fi
    rm -f "$FRONTEND_PID_FILE"
else
    echo "Frontend PID file not found."
fi

echo "All services stopped."