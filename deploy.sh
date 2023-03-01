#!/bin/bash

# Clone the repository from GitHub
git clone https://github.com/lc1993929/ChatGPTServer

# Change directory to the cloned repository
# shellcheck disable=SC2164
cd ChatGPTServer


go get ChatGPTServer
# Build the Go application
go build -o ChatServer

# Execute the compiled binary
./ChatServer -apiKey
