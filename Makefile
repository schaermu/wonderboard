BINARY_NAME=wonderboard
GOCOVER=go tool cover

build:
	npm run --prefix ui build
	GOARCH=amd64 GOOS=linux go build -o ./build/${BINARY_NAME} .

run:
	./build/${BINARY_NAME} serve

start: clean build run

clean:
	go clean
	go clean -testcache
	rm -rf ./build
	rm -rf ./ui/dist

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