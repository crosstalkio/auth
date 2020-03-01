PROTOS := $(wildcard *.proto)
PBGO := $(PROTOS:.proto=.pb.go)

REDISTORE := redistore

%.pb.go: %.proto
	protoc --go_out=. $<

all: $(PBGO) $(REDISTORE)
	go build .

$(REDISTORE):
	go build -o $@ ./cmd/redistore

clean:
	rm -f $(PBGO)
	rm -f $(REDISTORE)

.PHONY: all clean
