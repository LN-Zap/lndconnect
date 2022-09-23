default: install

dep:
	go get ./...

install: dep
	go install -v ./...

test:
	go test -v