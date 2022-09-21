BINARY_NAME=lndconnect
GOARCH=amd64
GOOS=linux

default: run

build:
	go build -o ${BINARY_NAME}

run: build
	./${BINARY_NAME} --lnddir .data -j -c

dep:
	go mod download

test: dep
	go test -v