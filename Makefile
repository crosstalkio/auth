PROTOS := $(wildcard *.proto)
PBGO := $(PROTOS:.proto=.pb.go)

CROSSAUTH := crossauth

all: $(PBGO) $(CROSSAUTH)
	go build .

%.pb.go: %.proto
	protoc --go_out=. $<

$(CROSSAUTH):
	go build -o $@ ./cmd/crossauth

clean:
	rm -f $(PBGO)
	rm -f $(CROSSAUTH)

.PHONY: all clean
