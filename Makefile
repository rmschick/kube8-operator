GOVERSION = $(shell go version | awk '{print $$3;}')
SOURCE_FILES?=./...

export PATH := ./bin:$(PATH)
export GOPRIVATE := github.com/FishtechCSOC
export CGO_ENABLED := 0

clean:
	rm -rf ./dist && rm -rf ./vendor
.PHONY: clean

upgrade:
	go get -t -u ./...
.PHONY: upgrade

vendor:
	go mod vendor
.PHONY: vendor

tidy:
	go mod tidy
.PHONY: tidy

lint:
	golangci-lint run --timeout=5m
.PHONY: lint

test:
	gotestsum -- -failfast -v -covermode count -timeout 5m $(SOURCE_FILES)
.PHONY: test

build:
	GOVERSION=$(GOVERSION) goreleaser release --snapshot --skip-publish --skip-sign --clean --debug
.PHONY: build

snapshot:
	GOVERSION=$(GOVERSION) goreleaser release --snapshot --clean --skip-sign --debug
.PHONY: snapshot

release:
	GOVERSION=$(GOVERSION) goreleaser release --clean --skip-sign --debug
.PHONY: release

docs:
	# Docs available at http://localhost:6060
	godoc -http=:6060
.PHONY: docs
