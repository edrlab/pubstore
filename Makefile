
ROOT_DIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

BUILD_DIR=$(ROOT_DIR)/build

rm=rm -rf

pubstore=cmd/pubstore/pubstore.go

.PHONY: all clean $(pubstore) test

all: $(pubstore)

clean:
	$(rm)  $(BUILD_DIR)

test:
	go test -coverpkg=./pkg/./... ./pkg/./...

$(pubstore):	
	GOPATH=$(BUILD_DIR) go install ./$@
	
