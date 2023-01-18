all: build test
lint: golangci-lint super-linter
pipeline: all lint

build:
	go build -v ./...

test:
	go test -v ./...

golangci-lint:
	golangci-lint run

super-linter:
	docker run \
		--rm \
		--volume "$(shell pwd):/work:z" \
		--env DEFAULT_WORKSPACE=/work \
		--env RUN_LOCAL=true \
		--env VALIDATE_GO=false \
		github/super-linter:v4.10.0 bash

.PHONY: all build golangci-lint lint pipeline super-linter test
