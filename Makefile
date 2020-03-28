PROTOS := $(wildcard *.proto)
PBGO := $(PROTOS:.proto=.pb.go)

PROTOGENGO := $(GOPATH)/bin/protoc-gen-go

CROSSAUTH := crossauth
GOFILES := go.mod $(wildcard *.go) $(wildcard */*.go)

all: $(PBGO) $(CROSSAUTH)
	go build .

$(PROTOGENGO):
	go install google.golang.org/protobuf/cmd/protoc-gen-go

%.pb.go: %.proto $(PROTOGENGO)
	protoc --go_out=. $<

$(CROSSAUTH): $(GOFILES)
	go build -o $@ ./cmd/crossauth

clean:
	rm -f $(PBGO)
	rm -f $(CROSSAUTH)

.PHONY: all clean
