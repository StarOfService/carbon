.PHONY: test
test:
	go test -v -cover -p 1 ./...

coverprofile:
	go test -v -coverprofile=coverprofile.out  -p 1  ./...
	go tool cover -html=coverprofile.out

build:
	go build -o carbon
