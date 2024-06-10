#!/bin/bash

# Get the list of supported architectures
platforms=$(go tool dist list)

# Directory to store the built test binaries
output_dir="./.bins"
mkdir -p "$output_dir"

# Loop through each platform and build the tests
for platform in $platforms; do
    echo "Building tests for platform: $platform"

    # Split platform into GOOS and GOARCH
    IFS='/' read -r GOOS GOARCH <<< "$platform"

    # Set the environment variables
    export GOOS=$GOOS
    export GOARCH=$GOARCH

    # Build the tests
    output_file="$output_dir/test_$GOOS_$GOARCH"
    if go test -c -o "$output_file"; then
        echo "Successfully built tests for platform: $platform"
    else
        echo "Failed to build tests for platform: $platform"
    fi
done

echo "Test binaries are stored in $output_dir"
