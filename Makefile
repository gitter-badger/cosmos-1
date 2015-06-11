MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
DIR := $(dir $(MAKEFILE_PATH))
GOPATH := $(DIR)vendor
COSMOS_VENDOR_PATH := $(GOPATH)/src/github.com/cosmos-io/cosmos

export DIR
export GOPATH
export COSMOS_VENDOR_PATH

default: build

build:
	@rm -rf $(COSMOS_VENDOR_PATH)

	@mkdir -p $(DIR)bin
	@mkdir -p $(COSMOS_VENDOR_PATH)

	@cp -r $(DIR)context $(COSMOS_VENDOR_PATH)/context
	@cp -r $(DIR)model $(COSMOS_VENDOR_PATH)/model
	@cp -r $(DIR)service $(COSMOS_VENDOR_PATH)/service
	@cp -r $(DIR)converter $(COSMOS_VENDOR_PATH)/converter
	@cp -r $(DIR)router $(COSMOS_VENDOR_PATH)/router
	@cp -r $(DIR)worker $(COSMOS_VENDOR_PATH)/worker
	@cp -r $(DIR)influxdb $(COSMOS_VENDOR_PATH)/influxdb

	go build -o $(DIR)bin/cosmos

run: build
	$(DIR)bin/cosmos

doc:
	godoc -http=:6060 -index

fmt:
	go fmt $(DIR)

lint:
	golint $(DIR)

test:
	go test $(DIR)

vet:
	go vet $(DIR)

clean:
	rm -rf $(DIR)bin