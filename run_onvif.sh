#!/bin/bash
# Script to run both the ONVIF backend and frontend on a Raspberry Pi
# Author: GitHub Copilot
# Created: May 6, 2025

# Get the IP address of the Pi
PI_IP=$(hostname -I | awk '{print $1}')
echo "Raspberry Pi IP address: $PI_IP"

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check for required dependencies
echo "Checking dependencies..."

# Check for Go
if ! command_exists go; then
    echo "Go is not installed. Installing Go..."
    sudo apt update
    sudo apt install -y golang-go
fi

# Check for Node.js and npm
if ! command_exists node; then
    echo "Node.js is not installed. Installing Node.js and npm..."
    sudo apt update
    sudo apt install -y nodejs npm
fi

# Directory setup
BASE_DIR=$(dirname "$(readlink -f "$0")")
BACKEND_DIR="$BASE_DIR/onvif_back"
FRONTEND_DIR="$BASE_DIR/onvif_front"

echo "Starting services..."
echo "Base directory: $BASE_DIR"
echo "Backend directory: $BACKEND_DIR"
echo "Frontend directory: $FRONTEND_DIR"

# Create a directory for logs
mkdir -p "$BASE_DIR/logs"

# Kill any existing processes
echo "Checking for existing processes..."
if [ -f "$BASE_DIR/logs/backend.pid" ]; then
    BACKEND_PID=$(cat "$BASE_DIR/logs/backend.pid")
    if ps -p $BACKEND_PID > /dev/null; then
        echo "Stopping existing backend process..."
        kill $BACKEND_PID
    fi
fi

if [ -f "$BASE_DIR/logs/frontend.pid" ]; then
    FRONTEND_PID=$(cat "$BASE_DIR/logs/frontend.pid")
    if ps -p $FRONTEND_PID > /dev/null; then
        echo "Stopping existing frontend process..."
        kill $FRONTEND_PID
    fi
fi

# Function to start the backend
start_backend() {
    echo "Starting ONVIF backend server..."
    cd "$BACKEND_DIR" || exit
    go run main.go > "$BASE_DIR/logs/backend.log" 2>&1 &
    BACKEND_PID=$!
    echo "Backend started with PID: $BACKEND_PID"
    echo $BACKEND_PID > "$BASE_DIR/logs/backend.pid"
    sleep 3  # Give it a bit more time to initialize
}

# Function to start the frontend
start_frontend() {
    echo "Starting ONVIF frontend server..."
    cd "$FRONTEND_DIR" || exit
    # Set environment variables for the frontend
    export BACKEND_HOST=$PI_IP
    export BACKEND_PORT=8090
    npm run dev -- --host 0.0.0.0 > "$BASE_DIR/logs/frontend.log" 2>&1 &
    FRONTEND_PID=$!
    echo "Frontend started with PID: $FRONTEND_PID"
    echo $FRONTEND_PID > "$BASE_DIR/logs/frontend.pid"
    sleep 5  # Give it more time to initialize
}

# Start both services
start_backend
start_frontend

# Check if services are actually running
sleep 2
if ! ps -p "$(cat $BASE_DIR/logs/backend.pid)" > /dev/null; then
    echo "WARNING: Backend process failed to start. Check logs at $BASE_DIR/logs/backend.log"
    echo "Last 10 lines of backend log:"
    tail -10 "$BASE_DIR/logs/backend.log"
fi

if ! ps -p "$(cat $BASE_DIR/logs/frontend.pid)" > /dev/null; then
    echo "WARNING: Frontend process failed to start. Check logs at $BASE_DIR/logs/frontend.log"
    echo "Last 10 lines of frontend log:"
    tail -10 "$BASE_DIR/logs/frontend.log"
fi

echo "=================================================="
echo "ONVIF application is now running!"
echo "You can access the web interface from your laptop at:"
echo "http://$PI_IP:5173"
echo ""
echo "To view the logs in another terminal:"
echo "  Backend: tail -f $BASE_DIR/logs/backend.log"
echo "  Frontend: tail -f $BASE_DIR/logs/frontend.log"
echo ""
echo "To stop the services, run the following command:"
echo "  $BASE_DIR/stop_onvif.sh"
echo "=================================================="

# Keep the script running to make it easy to stop with Ctrl+C
echo "Press Ctrl+C to stop both services..."
trap 'echo "Stopping services..."; kill $(cat $BASE_DIR/logs/backend.pid) $(cat $BASE_DIR/logs/frontend.pid) 2>/dev/null; echo "Services stopped."; exit 0' INT

# This sleep infinity command ensures the script keeps running until Ctrl+C is pressed
sleep infinity