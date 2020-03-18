PROTOS := $(wildcard *.proto)
PBGO := $(PROTOS:.proto=.pb.go)

CROSSAUTH := crossauth
GOFILES := go.mod $(wildcard *.go) $(wildcard */*.go)

all: $(PBGO) $(CROSSAUTH)
	go build .

%.pb.go: %.proto
	protoc --go_out=. $<

$(CROSSAUTH): $(GOFILES)
	go build -o $@ ./cmd/crossauth

clean:
	rm -f $(PBGO)
	rm -f $(CROSSAUTH)

.PHONY: all clean
