PROTOS := $(wildcard *.proto)
PBGO := $(PROTOS:.proto=.pb.go)

CROSSAUTH := crossauth
GOFILES := go.mod $(wildcard *.go) $(wildcard */*.go)

all: $(PBGO) $(CROSSAUTH)
	go build .

$(CROSSAUTH): $(GOFILES)
	go build -o $@ ./cmd/crossauth

clean: clean/proto
	rm -f $(CROSSAUTH)

test: # -count=1 disables cache
	go test -v -race -count=1 .

.PHONY: all clean test

include .make/lint.mk
include .make/proto.mk
