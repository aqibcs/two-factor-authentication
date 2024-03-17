#!/bin/bash

source ./scripts/env.sh

Build=false

# Parse command line arguments
while getopts "b" opt; do
    case $opt in
    b)
        Build=true
        ;;
    esac
done

# Execute the build script if -b flag is provided
[ "$Build" = true ] && ./scripts/build.sh

# Run the binary with the specified CSV file
"$BINARY"