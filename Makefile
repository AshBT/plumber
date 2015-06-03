.PHONY: build all clean test install
all: build test

bindata/bindata.go:
	mkdir -p bindata
	go-bindata -pkg="bindata" -o=$@ manager

build: bindata/bindata.go
	./scripts/do.sh build

install: bindata/bindata.go
	./scripts/do.sh install

test: build
	go test -cover ./...

clean:
	@rm bindata/* plumber
	@rmdir bindata
	go clean
