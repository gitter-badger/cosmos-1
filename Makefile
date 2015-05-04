GOPATH := ${PWD}/vendor:${GOPATH}
export GOPATH

default: build

build:
	rm -rf ./vendor/src/github.com/cosmos-io/cosmos
	mkdir -p ./vendor/src/github.com/cosmos-io/cosmos
	cp -r ./context ./vendor/src/github.com/cosmos-io/cosmos/context
	cp -r ./dao ./vendor/src/github.com/cosmos-io/cosmos/dao
	cp -r ./model ./vendor/src/github.com/cosmos-io/cosmos/model
	cp -r ./service ./vendor/src/github.com/cosmos-io/cosmos/service
	cp -r ./converter ./vendor/src/github.com/cosmos-io/cosmos/converter
	cp -r ./router ./vendor/src/github.com/cosmos-io/cosmos/router
	cp -r ./util ./vendor/src/github.com/cosmos-io/cosmos/util
	cp -r ./worker ./vendor/src/github.com/cosmos-io/cosmos/worker
	go build -v -o ./bin/cosmos

run: build
	rm -rf ./bin/telescope
	cp -r ./telescope ./bin/telescope
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