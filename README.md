[![CircleCI](https://circleci.com/gh/schaermu/wonderboard/tree/main.svg?style=shield)](https://circleci.com/gh/schaermu/wonderboard/tree/main)
[![codecov](https://codecov.io/gh/schaermu/wonderboard/branch/main/graph/badge.svg?token=QC1WL6JQTQ)](https://codecov.io/gh/schaermu/wonderboard)
# wonderboard
Wonderboard aims to provide homelab users with a zero-configuration, opinionated application dashboard for their docker-compose based service stacks. Service and URL discovery is being done using Docker and Traefik API's (if applicable).

# Setup
1. Clone the repository: `git clone https://github.com/schaermu/wonderboard.git`
2. Build and run the app: `make start`

# Build
Run `docker build .` to build a production-ready, [distroless](https://github.com/GoogleContainerTools/distroless)-based docker image. The latest version of the application is also published to the [Github Registry](https://github.com/schaermu/wonderboard/pkgs/container/wonderboard).

# Test
## Single run
* Run all tests: `make test`
* Only run Svelte tests: `make test-svelte`
* Only run Go tests: `make test-go`

## Watch mode
* Watch all tests: `make watch`
* Only watch Svelte tests: `make watch-svelte`
* Only watch Go tests: `make watch-go`

# Develop
If you want to have a seamless development experience, i recommend you install [Air](https://github.com/cosmtrek/air) to get "hot-reloading" in Go as well.
1) Change directory to `ui`.
2) Run `npm install`.
3) Run `npm run dev`.
3.1) *(optional)* Install `air` by running `go install github.com/cosmtrek/air@latest`.
4) In a new shell, run either `air` or `make run` within the root directory.

In case you have no docker containers running, you can boot up a small test stack using `docker compose -f testing/docker-compose.yml up -d`.