.PHONY: build all clean test install
all: build test

bindata/bindata.go:
	mkdir -p bindata
	go-bindata -pkg="bindata" -o=$@ manager

build: bindata/bindata.go
	sh scripts/do.sh build

install: bindata/bindata.go
	sh scripts/do.sh install

test: build
	go test -cover ./...

clean:
	@rm bindata/* plumb
	@rmdir bindata
	go clean
