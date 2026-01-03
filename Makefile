# Makefile for PubStore

pubstore=cmd/pubstore/pubstore.go

swag=~/go/bin/swag

.PHONY: all $(pubstore) test docs

all: $(pubstore)

build: $(pubstore)

$(pubstore):	
	go build -o $$GOPATH/bin/pubstore  ./cmd/pubstore/.

test:
	go test -coverpkg=./pkg/./... ./pkg/./...

docs:
	$(swag) init -g router.go -d pkg/api -o pkg/docs

