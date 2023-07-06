
ROOT_DIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

BUILD_DIR=$(ROOT_DIR)/build

rm=rm -rf

pubstore=cmd/pubstore/pubstore.go

swag=~/go/bin/swag

.PHONY: all clean $(pubstore) test docs run

all: $(pubstore)

clean:
	$(rm)  $(BUILD_DIR)

test:
	go test -coverpkg=./pkg/./... ./pkg/./...

build: $(pubstore)

docs:
	$(swag) init -g router.go -d pkg/api -o pkg/docs

$(pubstore):	
	GOPATH=$(BUILD_DIR) go install ./$@
	
run:
	./build/bin/pubstore
