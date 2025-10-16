#!/bin/bash

# Script to generate Go code from protobuf definitions

# Create proto directory if it doesn't exist
mkdir -p proto/media

# Copy media service proto from media service
cp ../media-service/proto/media_service.proto proto/media/

# Generate Go code
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/media/media_service.proto

echo "Proto generation complete!"
