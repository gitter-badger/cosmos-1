GOPATH := ${PWD}/vendor
export GOPATH

default: build

build:
	@rm -rf ./vendor/src/github.com/cosmos-io/cosmos
	@rm -rf ./bin/telescope
	@mkdir -p ./bin
	@mkdir -p ./vendor/src/github.com/cosmos-io/cosmos
	@cp -r ./context ./vendor/src/github.com/cosmos-io/cosmos/context
	@cp -r ./dao ./vendor/src/github.com/cosmos-io/cosmos/dao
	@cp -r ./model ./vendor/src/github.com/cosmos-io/cosmos/model
	@cp -r ./service ./vendor/src/github.com/cosmos-io/cosmos/service
	@cp -r ./converter ./vendor/src/github.com/cosmos-io/cosmos/converter
	@cp -r ./router ./vendor/src/github.com/cosmos-io/cosmos/router
	@cp -r ./worker ./vendor/src/github.com/cosmos-io/cosmos/worker
	@cp -r ./influxdb ./vendor/src/github.com/cosmos-io/cosmos/influxdb
	@cp -r ./telescope ./bin/telescope
	go build -o ./bin/cosmos

run: build
	./bin/cosmos

doc:
	godoc -http=:6060 -index

fmt:
	go fmt .

lint:
	golint .

test:
	go test .

vet:
	go vet .

clean:
	rm -rf ./bin