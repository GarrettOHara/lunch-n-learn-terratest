#!/bin/bash

# Build the Docker image
docker build -t amazonlinux-2-go .

# Create a temporary container from the built image
container_id=$(docker create amazonlinux-2-go)

# Copy the built executable from the container to the current directory
docker cp $container_id:/src/server $(pwd)/.

# Remove the temporary container
docker rm $container_id

# For inside the container, run commands:
# yum update -y ; yum install -y wget tar gzip git ; yum clean all
# wget https://golang.org/dl/go1.22.3.linux-amd64.tar.gz
# tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz
# rm go1.22.3.linux-amd64.tar.gz
# export PATH=$PATH:/usr/local/go/bin
# go version
# mkdir /src
# ls
# git --version
# git clone https://github.com/GarrettOHara/lunch-n-learn-terratest.git
# cd lunch-n-learn-terratest/src/server/
# ls
# rm Dockerfile chat.db go.mod go.sum
# ls
# go mod init server ; go mod tidy ; go build
# ./server &
# curl localhost:80/
