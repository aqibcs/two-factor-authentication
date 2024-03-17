#!/bin/bash

# Load environment variables from env.sh
source ./scripts/env.sh

# Check if the output directory exists, create it if not
if [ ! -d "$OUTPUT_DIRECTORY" ]; then
    mkdir -p "$OUTPUT_DIRECTORY"
    if [ $? -ne 0 ]; then
        echo "Error: Failed to create output directory: $OUTPUT_DIRECTORY"
        exit 1
    fi
fi

# Build the Go project
go build -o "$BINARY" "$SOURCE_FILE"
if [ $? -ne 0 ]; then
    echo "Error: Failed to build the project. Check the Go build logs for details."
    exit 1
fi

echo "Build project successfully. Binary file: $BINARY"