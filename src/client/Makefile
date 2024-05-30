BINARY_NAME=client

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows main.go

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm -f ${BINARY_NAME}-darwin
	rm -f ${BINARY_NAME}-linux
	rm -f ${BINARY_NAME}-windows

lint:
	golangci-lint run --enable-all

.PHONY: build run clean lint
