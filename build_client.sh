#!/bin/bash

# Build script for Go client

echo "Building Go client..."
cd client

# Build the client
go build -o model-client main.go

if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo "Client binary: client/model-client"
    echo ""
    echo "Usage examples:"
    echo "  ./client/model-client list"
    echo "  ./client/model-client download -name bert-base-uncased"
    echo "  ./client/model-client get -id bert-base-uncased"
    echo "  ./client/model-client delete -id bert-base-uncased"
    echo "  ./client/model-client health"
else
    echo "Build failed"
    exit 1
fi
