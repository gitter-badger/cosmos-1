GOPATH := ${PWD}/vendor:${GOPATH}
export GOPATH

default: build

build: vet
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