PROTOS := $(wildcard *.proto)
PBGO := $(PROTOS:.proto=.pb.go)

CROSSAUTH := crossauth
GOFILES := go.mod $(wildcard *.go) $(wildcard */*.go)

all: $(PBGO) $(CROSSAUTH)
	go build .

tidy:
	go mod tidy

$(CROSSAUTH): $(GOFILES)
	go build -o $@ ./cmd/crossauth

clean: clean/proto
	rm -f go.sum
	rm -f $(CROSSAUTH)

test: # -count=1 disables cache
	go test -v -race -count=1 .

.PHONY: all tidy clean test

include .make/lint.mk
include .make/proto.mk
