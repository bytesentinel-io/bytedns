#!/bin/bash

# Build the project
BUILD_DIR=build

rm -rf $BUILD_DIR

# Create the build directory if it doesn't exist
if [ ! -d "$BUILD_DIR" ]; then
    mkdir -p $BUILD_DIR
    mkdir -p $BUILD_DIR/linux/{amd64,386,arm64,arm}
    mkdir -p $BUILD_DIR/mac/{amd64,386,arm64,arm}
    mkdir -p $BUILD_DIR/windows/{amd64,386,arm64,arm}
fi

# Build go binary for all platforms
GOOS=linux GOARCH=amd64 go build -o $BUILD_DIR/linux/amd64/ ./src
# GOOS=linux GOARCH=386 go build -o $BUILD_DIR/linux/386/ ./...
GOOS=linux GOARCH=arm64 go build -o $BUILD_DIR/linux/arm64/ ./src
# GOOS=linux GOARCH=arm go build -o $BUILD_DIR/linux/arm/ ./...

GOOS=darwin GOARCH=amd64 go build -o $BUILD_DIR/mac/amd64/ ./src
# GOOS=darwin GOARCH=386 go build -o $BUILD_DIR/mac/386/ ./...
GOOS=darwin GOARCH=arm64 go build -o $BUILD_DIR/mac/arm64/ ./src
# GOOS=darwin GOARCH=arm go build -o $BUILD_DIR/mac/arm/ ./...

GOOS=windows GOARCH=amd64 go build -o $BUILD_DIR/windows/amd64/ ./src
# GOOS=windows GOARCH=386 go build -o $BUILD_DIR/windows/386/ ./...
GOOS=windows GOARCH=arm64 go build -o $BUILD_DIR/windows/arm64/ ./src
# GOOS=windows GOARCH=arm go build -o $BUILD_DIR/windows/arm/ ./...