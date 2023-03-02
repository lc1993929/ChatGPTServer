#!/bin/bash


# 设置要kill的进程名和端口号
PROC_NAME="ChatServer"
PORT_NUM="8080"

# 查找进程号
PID=$(netstat -nlp | grep :$PORT_NUM | awk '{print $7}' | awk -F"/" '{print $1}')

# 如果进程号存在，则kill进程
if [ -n "$PID" ]; then
    echo "Killing process $PROC_NAME (PID=$PID)..."
    kill $PID
else
    echo "No process found with name $PROC_NAME and port number $PORT_NUM."
fi

rm -rf ChatGPTServer

# Clone the repository from GitHub
git clone https://github.com/lc1993929/ChatGPTServer

# Change directory to the cloned repository
# shellcheck disable=SC2164
cd ChatGPTServer


go get ChatGPTServer
# Build the Go application
go build -o ChatServer

# Execute the compiled binary
./ChatServer -apiKey  &
