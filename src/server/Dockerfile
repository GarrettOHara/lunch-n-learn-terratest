# Use the official Amazon Linux 2 image as base
FROM amazonlinux:2

# Install required dependencies for installing Go
RUN yum update -y && \
    yum install -y wget tar gzip && \
    yum clean all

# Download and install Go 1.22.3
RUN wget https://golang.org/dl/go1.22.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz && \
    rm go1.22.3.linux-amd64.tar.gz

# Add Go binary to the system PATH
ENV PATH=$PATH:/usr/local/go/bin

# Create a directory 'src' in the container
RUN mkdir /src

# Set the working directory to 'src'
WORKDIR /src

# Copy the contents of the current directory on the host to the 'src' directory in the container
COPY . /src

RUN go mod init server && go mod tidy && go build

# Copy the built executable from the container to the host machine
CMD ["cp", "/src/server", "./build"]


# Print Go version to confirm installation
RUN go version

# Run commands:
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
