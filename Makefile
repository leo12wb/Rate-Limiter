include .env

MAIN_PACKAGE_PATH := ./cmd/server
BINARY_NAME := server

# ============================================================================ #
# HELPERS
# ============================================================================ #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ============================================================================ #
# Development
# ============================================================================ #

## test: run all tests
.PHONY: test
test:
		go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
		go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
		go tool cover -html=/tmp/coverage.out

## build: build the application
.PHONY: build
build:
	go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## run: run the  application
.PHONY: run
run: build
	/tmp/bin/$(BINARY_NAME)

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.51.0 \
			--build.cmd "make build" --build.bin "/tmp/bin/${BINARY_NAME}" --build.delay "100" \
			--build.exclude_dir "" \
			--build.include_ext "go" \
			--misc.clean_on_exit "true"