BINARY_NAME=wonderboard
GOCOVER=go tool cover

build:
	GOARCH=amd64 GOOS=linux go build -o ./build/${BINARY_NAME} .
build-all:
	npm run --prefix ui build
	GOARCH=amd64 GOOS=linux go build -o ./build/${BINARY_NAME} .

run:
	./build/${BINARY_NAME} serve

start: clean build run
start-all: clean build-all run

clean:
	go clean
	go clean -testcache
	rm -rf ./build

test: test-go test-svelte

test-go:
	gotestsum -f testname -- -tags=test -coverprofile=coverage.txt -race -covermode=atomic ./...

test-svelte:
	npm run --prefix ui test

watch:
	make -j2 watch-go watch-svelte

watch-go:
	gotestsum --watch -f testname -- -tags=test -coverprofile=coverage.txt -race -covermode=atomic ./...

watch-svelte:
	npm run --prefix ui test:watch

cover:
	gotestsum -f testname -- -tags=test ./... -coverprofile=coverage.out
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out -o coverage.html