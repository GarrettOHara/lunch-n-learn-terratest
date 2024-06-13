BINARY_NAME=client

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin-amd64
	GOARCH=arm64 GOOS=darwin go build -o ${BINARY_NAME}-darwin-arm64
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows.exe

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm -f ${BINARY_NAME}-darwin*
	rm -f ${BINARY_NAME}-linux
	rm -f ${BINARY_NAME}-windows.exe

lint:
	golangci-lint run --enable-all

.PHONY: build run clean lint
