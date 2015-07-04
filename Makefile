MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
DIR := $(dir $(MAKEFILE_PATH))
GOPATH := $(DIR)vendor
COSMOS_VENDOR_PATH := $(GOPATH)/src/github.com/cosmos-io/cosmos

export DIR
export GOPATH
export COSMOS_VENDOR_PATH
export CGO_ENABLED=0

default: build

build:
	@rm -rf $(COSMOS_VENDOR_PATH)
	@rm -rf $(DIR)bin/telescope

	@mkdir -p $(DIR)bin
	@mkdir -p $(COSMOS_VENDOR_PATH)

	@cp -r $(DIR)context $(COSMOS_VENDOR_PATH)/context
	@cp -r $(DIR)model $(COSMOS_VENDOR_PATH)/model
	@cp -r $(DIR)route $(COSMOS_VENDOR_PATH)/route
	@cp -r $(DIR)influxdb $(COSMOS_VENDOR_PATH)/influxdb
	@cp -r $(DIR)telescope $(DIR)bin

	go build -a -ldflags '-s' -o $(DIR)bin/cosmos

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