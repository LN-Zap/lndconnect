default: install

dep:
	go get ./...

install: dep
	go install -v ./...
