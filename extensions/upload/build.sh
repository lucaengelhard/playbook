#!/bin/bash

# https://hackernoon.com/one-go-codebase-five-platforms-my-cross-platform-build-journey

echo "Building upload extension"

GOOS=windows GOARCH=amd64 go build -o "bin/windows/playbook-upload-amd64.exe"

GOOS=linux GOARCH=amd64 go build -o "bin/linux/playbook-upload-amd64"
#GOOS=linux GOARCH=arm64 go build -o "bin/linux/playbook-upload-arm64"

GOOS=darwin GOARCH=amd64 go build -o "bin/mac/playbook-upload-amd64"
#GOOS=darwin GOARCH=arm64 go build -o "bin/mac/playbook-upload-arm64"

echo "Done :)"
