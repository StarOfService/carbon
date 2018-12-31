.PHONY: all codestyle minikube test coverprofile build build-all-os

all: build

codestyle:
	gometalinter  ./cmd/  ./pkg/...
	# TODO: revie checks

minikube:
	@minikube status > /dev/null || minikube start

test: minikube
	@set -e
	@minikube status > /dev/null
	go test -count=1 -p 1 `go list ./... | grep -v minikube` github.com/starofservice/carbon/pkg/minikube

coverprofile: minikube
	go test -v -coverprofile=coverprofile.out -p 1 `go list ./... | grep -v minikube` github.com/starofservice/carbon/pkg/minikube
	go tool cover -html=coverprofile.out

build:
	go build -o carbon

build-all-os:
	env GOOS=windows GOARCH=amd64 go build -o carbon-windows-amd64.exe
	env GOOS=linux GOARCH=amd64 go build -o carbon-linux-amd64
	env GOOS=darwin GOARCH=amd64 go build -o carbon-darwin-amd64
