PROTOS := $(wildcard *.proto)
PBGO := $(PROTOS:.proto=.pb.go)

%.pb.go: %.proto
	protoc --go_out=. $<

all: $(PBGO)
	go build .

clean:
	rm -f $(PBGO)

.PHONY: all clean
