#!/bin/bash

# Start server script for Model Management System

# Set default models mount point if not set
if [ -z "$MODELS_MOUNT_POINT" ]; then
    export MODELS_MOUNT_POINT="./models"
    echo "Using default models mount point: $MODELS_MOUNT_POINT"
fi

# Check if Python dependencies are installed
if ! python3 -c "import flask" 2>/dev/null; then
    echo "Installing Python dependencies..."
    pip install -r requirements.txt
fi

# Create models directory if it doesn't exist
mkdir -p "$MODELS_MOUNT_POINT"

# Start the Flask server
echo "Starting Flask server..."
echo "Server will be available at http://localhost:5000"
echo "User view: http://localhost:5000/"
echo "Admin panel: http://localhost:5000/admin"
echo ""
cd server && python3 app.py
