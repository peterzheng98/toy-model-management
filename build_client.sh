#!/bin/bash

# Build script for Go client

echo "Building Go client..."
cd client/go

# Build the client
go build -o model-client main.go

if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo "Client binary: client/go/model-client"
    echo ""
    echo "Usage examples:"
    echo "  ./client/go/model-client list"
    echo "  ./client/go/model-client download -name bert-base-uncased"
    echo "  ./client/go/model-client get -id bert-base-uncased"
    echo "  ./client/go/model-client delete -id bert-base-uncased"
    echo "  ./client/go/model-client health"
    echo ""
    echo "Python client is also available:"
    echo "  python client/python/model_client.py list"
else
    echo "Build failed"
    exit 1
fi
