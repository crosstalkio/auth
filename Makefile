PROTOS := $(wildcard *.proto)
PBGO := $(PROTOS:.proto=.pb.go)

CROSSAUTH := crossauth
GOFILES := go.mod $(wildcard *.go) $(wildcard */*.go) $(wildcard */*/*.go)

all: $(PBGO) $(CROSSAUTH)
	go build .

include .make/golangci-lint.mk
include .make/protoc.mk
include .make/protoc-gen-go.mk

$(CROSSAUTH): $(GOFILES)
	go build -o $@ ./cmd/crossauth

tidy:
	go mod tidy

lint: $(GOLANGCI_LINT)
	$(realpath $(GOLANGCI_LINT)) run

clean: clean/golangci-lint clean/protoc clean/protoc-gen-go
	rm -f go.sum
	rm -f $(PBGO)
	rm -f $(CROSSAUTH)

test: # -count=1 disables cache
	go test -v -race -count=1 .

.PHONY: all tidy lint clean test
