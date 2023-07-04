
ROOT_DIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

BUILD_DIR=$(ROOT_DIR)/build

rm=rm -rf

pubstore=cmd/pubstore/pubstore.go

.PHONY: all clean $(pubstore)

all: $(pubstore)

clean:
	$(rm)  $(BUILD_DIR)

$(pubstore):	
	GOPATH=$(BUILD_DIR) go install ./$@
	
